// Package gdbinit 把项目统一数据库配置（DB_DSN）接入 GoFrame ORM 数据源，
// 使 dao 层与业务层统一指向 MySQL 数据库。
package gdbinit

import (
	"fmt"
	"strings"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
)

// Setup 注册名为 default 的 ORM 数据源（MySQL），dsn 为 go-sql-driver 格式连接串。
func Setup(dsn string) error {
	node := &gdb.ConfigNode{Type: "mysql"}
	user, pass, addr, name, err := parseMySQLDSN(dsn)
	if err != nil {
		// 非 go-sql-driver 格式（如 "user:pass@tcp(host:port)/db" 之外）退回 Link 方式
		node.Link = dsn
	} else {
		node.User, node.Pass, node.Name = user, pass, name
		node.Host, node.Port = splitHostPort(addr, "3306")
		node.Extra = "parseTime=true&loc=Local&charset=utf8mb4"
	}
	if err := gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{*node},
	}); err != nil {
		return fmt.Errorf("gdbinit: 注册数据源失败: %w", err)
	}
	return nil
}

// parseMySQLDSN 解析 go-sql-driver 格式 DSN：user:pass@tcp(host:port)/dbname?params
func parseMySQLDSN(dsn string) (user, pass, addr, name string, err error) {
	at := strings.LastIndex(dsn, "@")
	slash := strings.LastIndex(dsn, "/")
	if at < 0 || slash < at {
		return "", "", "", "", fmt.Errorf("invalid dsn")
	}
	cred := dsn[:at]
	if c := strings.Index(cred, ":"); c >= 0 {
		user, pass = cred[:c], cred[c+1:]
	} else {
		user = cred
	}
	name = dsn[slash+1:]
	if q := strings.Index(name, "?"); q >= 0 {
		name = name[:q]
	}
	rest := dsn[at+1 : slash]
	rest = strings.TrimPrefix(rest, "tcp(")
	rest = strings.TrimSuffix(rest, ")")
	addr = rest
	return user, pass, addr, name, nil
}

func splitHostPort(addr, defPort string) (string, string) {
	if c := strings.LastIndex(addr, ":"); c >= 0 {
		return addr[:c], addr[c+1:]
	}
	return addr, defPort
}
