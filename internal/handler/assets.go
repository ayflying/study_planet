package handler

import (
	"embed"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

//go:embed assets
var assetsFS embed.FS

// Logo 对外提供品牌 Logo（GET /assets/logo.svg）。
func (s *Store) Logo(r *ghttp.Request) {
	b, err := assetsFS.ReadFile("assets/logo.svg")
	if err != nil {
		r.Response.WriteStatus(http.StatusNotFound, "logo missing")
		return
	}
	r.Response.Header().Set("Content-Type", "image/svg+xml")
	r.Response.Header().Set("Cache-Control", "public, max-age=86400")
	r.Response.Write(b)
}

// Index 提供多邻国式学习主页（GET / 及 /app）。
// 内容为内嵌单文件 HTML，零外链，移动端优先。
func (s *Store) Index(r *ghttp.Request) {
	b, err := assetsFS.ReadFile("assets/app.html")
	if err != nil {
		r.Response.WriteStatus(http.StatusNotFound, "app missing")
		return
	}
	html := string(b)
	// 按请求推导 API 基地址，避免硬编码域名/端口
	base := ""
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		scheme := r.Header.Get("X-Forwarded-Proto")
		if scheme == "" {
			scheme = "http"
		}
		base = scheme + "://" + h
	}
	html = strings.ReplaceAll(html, "__API_BASE__", base+"/api")
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Write(html)
}
