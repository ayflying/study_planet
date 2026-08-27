package db

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动（CGO-free，便于交叉编译与日后切 Postgres/MySQL）
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate 使用 golang-migrate 跑全部 up 迁移（幂等，无变更时返回 ErrNoChange 也视为成功）。
// dsn 为 SQLite 文件路径，例如 "data/grade5.db"。
func Migrate(dsn string) error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	// 复用同一个 modernc 连接执行迁移，避免 DSN 解析歧义
	conn, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	driver, err := sqlite.WithInstance(conn.DB, &sqlite.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", src, "sqlite", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// Open 打开 SQLite 连接池。SQLite 为单写者，限制单连接避免锁竞争。
func Open(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
