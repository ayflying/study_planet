package middleware

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
)

// CORS 开发态放行跨域，并正确处理 OPTIONS 预检。
func CORS(r *ghttp.Request) {
	r.Response.Header().Set("Access-Control-Allow-Origin", "*")
	r.Response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	if r.Method == "OPTIONS" {
		r.Response.WriteStatusExit(http.StatusNoContent)
		return
	}
	r.Middleware.Next()
}

// ParentAuth 校验 Authorization: Bearer <token>，失败直接 401 中断。
// 401 属于 HTTP 鉴权语义（前端据此清理本地 token），因此保留状态码而非业务 code。
func ParentAuth(secret string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		auth := r.Header.Get("Authorization")
		tok := strings.TrimPrefix(auth, "Bearer ")
		if tok == "" {
			r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{"code": http.StatusUnauthorized, "message": "未登录", "data": nil})
			return
		}
		if _, err := jwt.Parse(tok, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}); err != nil {
			r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{"code": http.StatusUnauthorized, "message": "登录已失效", "data": nil})
			return
		}
		r.Middleware.Next()
	}
}
