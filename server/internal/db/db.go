package db

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path"
	"regexp"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*/*.sql
var migrationsFS embed.FS

type migrationFile struct {
	version int
	name    string
	upSQL   string
}

func loadMigrations(driver string) ([]migrationFile, error) {
	dir := path.Join("migrations", driver)
	entries, err := fs.ReadDir(migrationsFS, dir)
	if err != nil {
		return nil, err
	}
	var ups []migrationFile
	seen := map[int]string{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}
		content, err := migrationsFS.ReadFile(path.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		var version int
		if _, err := fmt.Sscanf(e.Name(), "%06d", &version); err != nil {
			return nil, fmt.Errorf("解析迁移版本失败 %s: %w", e.Name(), err)
		}
		// 同版本号只允许一个文件：重复会导致其中一个永远不会被执行（历史踩坑：000008 双文件）。
		if prev, dup := seen[version]; dup {
			return nil, fmt.Errorf("迁移版本 %06d 重复：%s 与 %s", version, prev, e.Name())
		}
		seen[version] = e.Name()
		ups = append(ups, migrationFile{version: version, name: e.Name(), upSQL: string(content)})
	}
	sort.Slice(ups, func(i, j int) bool { return ups[i].version < ups[j].version })
	return ups, nil
}

func splitStatements(sqlText string) []string {
	// 归一化行尾：git checkout 在 Windows 下会把迁移脚本转成 CRLF，
	// 残留的 \r 会被 MySQL 当作语句的一部分导致语法错误。
	sqlText = strings.ReplaceAll(sqlText, "\r\n", "\n")
	lines := strings.Split(sqlText, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		cleaned = append(cleaned, line)
	}
	text := strings.Join(cleaned, "\n")
	parts := strings.Split(text, ";\n")
	var stmts []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimSuffix(p, ";")
		if p != "" {
			stmts = append(stmts, p)
		}
	}
	return stmts
}

func ensureMigrationsTable(db *sqlx.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL DEFAULT '', applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	return err
}

func appliedVersions(db *sqlx.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	applied := map[int]bool{}
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, rows.Err()
}

// existingColumns 返回指定表已存在的列名集合（小写）。用于让迁移脚本对
// "列已存在" 的场景幂等（MySQL 8 不支持 ADD COLUMN IF NOT EXISTS）。
func existingColumns(db *sqlx.DB, table string) (map[string]bool, error) {
	rows, err := db.Query(
		"SELECT column_name FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ?", table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols := map[string]bool{}
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		cols[strings.ToLower(c)] = true
	}
	return cols, rows.Err()
}

// statementCovered 判断单条语句是否因目标已存在而可以跳过。
// 当前只识别 "ALTER TABLE <t> ADD COLUMN <c> ..." 这一形态。
func statementCovered(db *sqlx.DB, stmt string) (bool, error) {
	m := alterAddColRe.FindStringSubmatch(stmt)
	if m == nil {
		return false, nil
	}
	cols, err := existingColumns(db, m[1])
	if err != nil {
		return false, err
	}
	return cols[strings.ToLower(m[2])], nil
}

// alterAddColRe 匹配 "ALTER TABLE <t> ADD COLUMN <c> ..."（容忍反引号/双引号包名
// 与 MariaDB 风格的 IF NOT EXISTS 残留写法，那些写法在 MySQL 8 本身会语法报错）。
var alterAddColRe = regexp.MustCompile(
	`(?is)^ALTER\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?[` + "`" + `"]?([a-zA-Z0-9_]+)[` + "`" + `"]?\s+` +
		`ADD\s+COLUMN\s+(?:IF\s+NOT\s+EXISTS\s+)?[` + "`" + `"]?([a-zA-Z0-9_]+)`)

// Migrate 启动时自动升级数据库结构（仅 MySQL）：空库全量建表，已有库按版本表增量执行；无变更时跳过。
func Migrate(dsn string) error {
	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	if err := ensureMigrationsTable(conn); err != nil {
		return fmt.Errorf("初始化迁移版本表失败: %w", err)
	}
	applied, err := appliedVersions(conn)
	if err != nil {
		return fmt.Errorf("读取迁移版本失败: %w", err)
	}
	migrations, err := loadMigrations("mysql")
	if err != nil {
		return fmt.Errorf("加载迁移脚本失败: %w", err)
	}
	var maxApplied int
	for v := range applied {
		if v > maxApplied {
			maxApplied = v
		}
	}
	var appliedCount int
	for _, m := range migrations {
		if applied[m.version] {
			continue
		}
		// 版本落后于已应用的最大版本：数据库结构与脚本不一致，拒绝静默执行。
		if m.version < maxApplied {
			return fmt.Errorf("数据库迁移版本不一致：已应用版本 %d 大于待执行版本 %d，请检查数据库", maxApplied, m.version)
		}
		for _, stmt := range splitStatements(m.upSQL) {
			// 兼容历史库结构与脚本不完全一致的情况：目标列已存在则跳过该语句，
			// 避免（且仅避免）"Duplicate column" 类可安全跳过的失败。
			if covered, err := statementCovered(conn, stmt); err != nil {
				return fmt.Errorf("检查迁移 %s 前置条件失败: %w", m.name, err)
			} else if covered {
				log.Printf("迁移 %s：目标已存在，跳过语句：%s", m.name, firstLine(stmt))
				continue
			}
			if _, err := conn.Exec(stmt); err != nil {
				return fmt.Errorf("执行迁移 %s 失败: %w", m.name, err)
			}
		}
		if _, err := conn.Exec("INSERT INTO schema_migrations(version,name) VALUES(?,?)", m.version, m.name); err != nil {
			return fmt.Errorf("记录迁移版本 %d 失败: %w", m.version, err)
		}
		appliedCount++
		maxApplied = m.version
	}
	log.Printf("数据库迁移完成：已应用 %d 个新迁移，当前最新版本 %06d", appliedCount, maxApplied)
	return nil
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
