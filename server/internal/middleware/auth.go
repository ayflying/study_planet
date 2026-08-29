package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
)

// ctxKeyParentID 请求上下文中家长 id 的键（logic 层经 ctxParentID 读取）。
type ctxKeyParentID struct{}

// ParentIDOf 从请求上下文取当前登录家长 id（未携带时为 0）。
func ParentIDOf(ctx context.Context) int {
	if v := ctx.Value(ctxKeyParentID{}); v != nil {
		if n, ok := v.(int); ok {
			return n
		}
	}
	return 0
}

// injectParentID 解析 JWT claims 并把 parent_id 写入请求上下文（解析失败静默跳过）。
func injectParentID(r *ghttp.Request, secret string) bool {
	auth := r.Header.Get("Authorization")
	tok := strings.TrimPrefix(auth, "Bearer ")
	if tok == "" {
		return false
	}
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(tok, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}); err != nil {
		return false
	}
	if raw, ok := claims["parent_id"]; ok {
		var pid int
		switch v := raw.(type) {
		case float64:
			pid = int(v)
		case int:
			pid = v
		}
		r.SetCtx(context.WithValue(r.Context(), ctxKeyParentID{}, pid))
		return true
	}
	return false
}

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
// 鉴权通过后把 claims 中的 parent_id 注入请求上下文，供业务层做数据归属隔离。
func ParentAuth(secret string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		auth := r.Header.Get("Authorization")
		tok := strings.TrimPrefix(auth, "Bearer ")
		if tok == "" {
			r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{"code": http.StatusUnauthorized, "message": "未登录", "data": nil})
			return
		}
		if !injectParentID(r, secret) {
			r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{"code": http.StatusUnauthorized, "message": "登录已失效", "data": nil})
			return
		}
		r.Middleware.Next()
	}
}

// ParentAuthOptional 可选鉴权：带有效 token 则注入家长身份，否则匿名放行。
// 用于「公开但登录后按身份过滤」的接口（如 /api/students 学生端切换器）。
func ParentAuthOptional(secret string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		injectParentID(r, secret)
		r.Middleware.Next()
	}
}
