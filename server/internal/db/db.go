package db

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	// SQLite 驱动由 GF contrib（github.com/glebarez/go-sqlite）提供，
	// 两者都注册 "sqlite" 驱动名会 panic（Register called twice），这里只引入一处。
	_ "github.com/glebarez/go-sqlite"
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
		ups = append(ups, migrationFile{version: version, name: e.Name(), upSQL: string(content)})
	}
	sort.Slice(ups, func(i, j int) bool { return ups[i].version < ups[j].version })
	return ups, nil
}

func splitStatements(sqlText string) []string {
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

func ensureMigrationsTable(db *sqlx.DB, driver string) error {
	switch driver {
	case "mysql":
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT NOT NULL PRIMARY KEY, name VARCHAR(255) NOT NULL DEFAULT '', applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
		return err
	default:
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER NOT NULL PRIMARY KEY, name TEXT NOT NULL DEFAULT '', applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
		return err
	}
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

// convertLegacyMigrations 兼容旧 golang-migrate 版本表：旧格式只有 version/dirty/ts 列、
// 只存一行最新版本；新格式按“每个已应用版本一行”记录（version+name）。
// 检测到旧格式时，重建为带 name 列的新表并补记 1..最新版本。
func convertLegacyMigrations(db *sqlx.DB) error {
	var row struct {
		Version int
		Dirty   bool
	}
	if err := db.Get(&row, "SELECT version, dirty FROM schema_migrations LIMIT 1"); err != nil {
		return nil // 没有 dirty 列，已是新格式
	}
	if row.Dirty {
		return fmt.Errorf("数据库存在未完成的迁移（dirty 版本 %d），请先人工修复 schema_migrations 表", row.Version)
	}
	if _, err := db.Exec("ALTER TABLE schema_migrations RENAME TO schema_migrations_legacy"); err != nil {
		return fmt.Errorf("备份旧迁移版本表失败: %w", err)
	}
	if err := ensureMigrationsTable(db, "sqlite"); err != nil {
		return fmt.Errorf("重建迁移版本表失败: %w", err)
	}
	for v := 1; v <= row.Version; v++ {
		if _, err := db.Exec("INSERT INTO schema_migrations(version,name) VALUES(?,?)", v, fmt.Sprintf("legacy-%06d", v)); err != nil {
			return fmt.Errorf("转换旧迁移版本记录失败（版本 %d）: %w", v, err)
		}
	}
	if _, err := db.Exec("DROP TABLE schema_migrations_legacy"); err != nil {
		return fmt.Errorf("清理旧迁移版本表失败: %w", err)
	}
	return nil
}

// Migrate 启动时自动升级数据库结构：空库全量建表，已有库按版本表增量执行；无变更时跳过。
func Migrate(driver, dsn string) error {
	driver = strings.ToLower(driver)
	if driver == "sqlite3" {
		driver = "sqlite"
	}
	conn, err := sqlx.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	if err := ensureMigrationsTable(conn, driver); err != nil {
		return fmt.Errorf("初始化迁移版本表失败: %w", err)
	}
	// 兼容旧 golang-migrate 单行版本表（SQLite 老库遗留）
	if driver != "mysql" {
		if err := convertLegacyMigrations(conn); err != nil {
			return err
		}
	}
	applied, err := appliedVersions(conn)
	if err != nil {
		return fmt.Errorf("读取迁移版本失败: %w", err)
	}
	migrations, err := loadMigrations(driver)
	if err != nil {
		return fmt.Errorf("加载迁移脚本失败: %w", err)
	}
	var maxApplied int
	for v := range applied {
		if v > maxApplied {
			maxApplied = v
		}
	}
	for _, m := range migrations {
		if applied[m.version] {
			continue
		}
		// 版本落后于已应用的最大版本：数据库结构与脚本不一致，拒绝静默执行。
		if m.version < maxApplied {
			return fmt.Errorf("数据库迁移版本不一致：已应用版本 %d 大于待执行版本 %d，请检查数据库", maxApplied, m.version)
		}
		for _, stmt := range splitStatements(m.upSQL) {
			if _, err := conn.Exec(stmt); err != nil {
				return fmt.Errorf("执行迁移 %s 失败: %w", m.name, err)
			}
		}
		if _, err := conn.Exec("INSERT INTO schema_migrations(version,name) VALUES(?,?)", m.version, m.name); err != nil {
			return fmt.Errorf("记录迁移版本 %d 失败: %w", m.version, err)
		}
	}
	return nil
}

// Open 打开数据库连接池。
func Open(driver, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(strings.ToLower(driver), dsn)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(driver, "sqlite") || strings.EqualFold(driver, "sqlite3") {
		db.SetMaxOpenConns(1)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
