package handler

import (
	"embed"
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
)

//go:embed assets/logo.svg
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
