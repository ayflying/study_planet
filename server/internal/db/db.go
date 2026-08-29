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
	// SQLite 驱动（纯 Go 实现，无需 CGO）由 GF contrib 同源提供；
	// 只在此处注册一次驱动名 "sqlite"，GF contrib/drivers/sqlite 复用同一驱动。
	_ "github.com/glebarez/go-sqlite"
)

//go:embed migrations/*/*.sql
var migrationsFS embed.FS

type migrationFile struct {
	version int
	name    string
	upSQL   string
}

// normDriver 规范驱动名（sqlite3 → sqlite）。
func normDriver(driver string) string {
	driver = strings.ToLower(driver)
	if driver == "sqlite3" {
		return "sqlite"
	}
	return driver
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
	if driver == "sqlite" {
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER NOT NULL PRIMARY KEY, name TEXT NOT NULL DEFAULT '', applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)")
		return err
	}
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

// Migrate 启动时自动升级数据库结构：空库全量建表，已有库按版本表增量执行；无变更时跳过。
// 支持 mysql 与 sqlite（本地开发免配置）双驱动。
func Migrate(driver, dsn string) error {
	driver = normDriver(driver)
	if driver != "mysql" && driver != "sqlite" {
		return fmt.Errorf("db: 不支持的数据库驱动 %s（仅支持 mysql/sqlite）", driver)
	}
	conn, err := sqlx.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer conn.Close()
	// SQLite 并发写限制：单连接避免 database is locked
	if driver == "sqlite" {
		conn.SetMaxOpenConns(1)
	}
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}
	if err := ensureMigrationsTable(conn, driver); err != nil {
		return fmt.Errorf("初始化迁移版本表失败: %w", err)
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
