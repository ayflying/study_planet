package config

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// Config 是服务端运行配置，优先读 manifest/config/config.yaml，环境变量可覆盖（容器/生产友好）。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Parent   ParentConfig
	Seed     SeedConfig
	Casdoor  CasdoorConfig
}

type ServerConfig struct {
	Port int
	CORS string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type ParentConfig struct {
	Pin       string
	JWTSecret string
}

type SeedConfig struct {
	Enabled bool
}

// CasdoorConfig Casdoor SSO（OIDC 授权码流程）。三项核心配齐即视为启用。
// 回调地址不配置，始终按用户访问的地址实时推导（见 handler.RedirectURIOf）。
type CasdoorConfig struct {
	Endpoint     string // 例如 https://casdoor.example.com
	ClientID     string
	ClientSecret string
	Organization string // 组织名，用于构造授权地址；留空则用应用默认
	Application  string // 应用名，回调 state 校验用
}

// Enabled 判断 Casdoor 是否已配置启用。
func (c *CasdoorConfig) Enabled() bool {
	return c.Endpoint != "" && c.ClientID != "" && c.ClientSecret != ""
}

// AuthURL 构造授权端点地址（Casdoor 标准 OIDC 路径）。redirectUri 由调用方按请求传入。
func (c *CasdoorConfig) AuthURL(redirectUri string) string {
	org := c.Organization
	if org == "" {
		org = "built-in"
	}
	return strings.TrimRight(c.Endpoint, "/") + "/login/oauth/authorize" +
		"?response_type=code" +
		"&client_id=" + url.QueryEscape(c.ClientID) +
		"&redirect_uri=" + url.QueryEscape(redirectUri) +
		"&scope=" + url.QueryEscape("read") +
		"&state=" + url.QueryEscape(c.Application) +
		"&org=" + url.QueryEscape(org)
}

// TokenURL 换取 access_token 的端点。
func (c *CasdoorConfig) TokenURL() string {
	return strings.TrimRight(c.Endpoint, "/") + "/api/login/oauth/access_token"
}

// UserURL 拉取用户信息的端点。
func (c *CasdoorConfig) UserURL() string {
	return strings.TrimRight(c.Endpoint, "/") + "/api/userinfo"
}

func getInt(c *gcfg.Config, ctx context.Context, key string, def int) int {
	if v, err := c.Get(ctx, key, def); err == nil {
		return v.Int()
	}
	return def
}

func getStr(c *gcfg.Config, ctx context.Context, key, def string) string {
	if v, err := c.Get(ctx, key, def); err == nil {
		return v.String()
	}
	return def
}

func getBool(c *gcfg.Config, ctx context.Context, key string, def bool) bool {
	if v, err := c.Get(ctx, key, def); err == nil {
		return v.Bool()
	}
	return def
}

// Load 读取配置；环境变量优先级高于文件：
//
//	SERVER_PORT / DB_DRIVER（仅支持 mysql） / DB_DSN / PARENT_PIN / JWT_SECRET
func Load() *Config {
	c := g.Cfg()
	ctx := context.Background()

	cfg := &Config{}
	cfg.Server.Port = getInt(c, ctx, "server.port", 8080)
	cfg.Server.CORS = getStr(c, ctx, "server.cors", "*")
	cfg.Database.Driver = getStr(c, ctx, "database.driver", "mysql")
	cfg.Database.DSN = getStr(c, ctx, "database.dsn", "")
	cfg.Parent.Pin = getStr(c, ctx, "parent.pin", "1234")
	cfg.Parent.JWTSecret = getStr(c, ctx, "parent.jwtSecret", "change-me-in-prod")
	cfg.Seed.Enabled = getBool(c, ctx, "seed.enabled", true)
	cfg.Casdoor.Endpoint = getStr(c, ctx, "casdoor.endpoint", "")
	cfg.Casdoor.ClientID = getStr(c, ctx, "casdoor.clientId", "")
	cfg.Casdoor.ClientSecret = getStr(c, ctx, "casdoor.clientSecret", "")
	cfg.Casdoor.Organization = getStr(c, ctx, "casdoor.organization", "")
	cfg.Casdoor.Application = getStr(c, ctx, "casdoor.application", "studyplanet")

	if v := os.Getenv("SERVER_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = n
		}
	}
	if v := os.Getenv("DB_DRIVER"); v != "" {
		cfg.Database.Driver = v
	}
	if v := os.Getenv("DB_DSN"); v != "" {
		cfg.Database.DSN = v
	}
	if v := os.Getenv("PARENT_PIN"); v != "" {
		cfg.Parent.Pin = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Parent.JWTSecret = v
	}
	// Casdoor：环境变量优先
	if v := os.Getenv("CASDOOR_ENDPOINT"); v != "" {
		cfg.Casdoor.Endpoint = v
	}
	if v := os.Getenv("CASDOOR_CLIENT_ID"); v != "" {
		cfg.Casdoor.ClientID = v
	}
	if v := os.Getenv("CASDOOR_CLIENT_SECRET"); v != "" {
		cfg.Casdoor.ClientSecret = v
	}
	if v := os.Getenv("CASDOOR_ORG_NAME"); v != "" {
		cfg.Casdoor.Organization = v
	}
	if v := os.Getenv("CASDOOR_APP_NAME"); v != "" {
		cfg.Casdoor.Application = v
	}
	return cfg
}
