package studyplanet

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// RedirectURIOf 按本次请求的访问地址实时推导回调地址（不写死、不落配置）。
// 依次识别 X-Forwarded-Proto / X-Forwarded-Host（反代场景），否则用请求自身 scheme+host。
// 这属于 HTTP 边缘语义，因此保留在控制器层。
func RedirectURIOf(r *ghttp.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	} else if p := r.Header.Get("X-Forwarded-Proto"); p != "" {
		scheme = p
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.GetHost()
	}
	return fmt.Sprintf("%s://%s/api/parent/casdoor/callback", scheme, host)
}

// CasdoorLogin 跳转 Casdoor 授权页（302）。
// logic 只负责计算授权页地址，跳转这一 HTTP 语义由控制器完成。
func (c *ControllerV1) CasdoorLogin(ctx context.Context, req *v1.CasdoorLoginReq) (res *v1.CasdoorLoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("请求上下文缺失")
	}
	req.RedirectURI = RedirectURIOf(r)
	res, err = svc.StudyPlanet().CasdoorLogin(ctx, req)
	if err != nil || res == nil || res.Location == "" {
		return res, err
	}
	r.Response.Header().Set("Location", res.Location)
	r.Response.WriteStatusExit(http.StatusFound)
	return
}

// CasdoorCallback Casdoor 授权码回调：logic 完成换 token 与落库，
// 这里以 text/html 输出回跳页（写 localStorage 后跳首页）。
func (c *ControllerV1) CasdoorCallback(ctx context.Context, req *v1.CasdoorCallbackReq) (res *v1.CasdoorCallbackRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("请求上下文缺失")
	}
	req.RedirectURI = RedirectURIOf(r)
	req.IsTLS = r.TLS != nil
	if req.Code == "" {
		req.Code = r.GetQuery("code").String()
	}
	res, err = svc.StudyPlanet().CasdoorCallback(ctx, req)
	if err != nil || res == nil || res.HTML == "" {
		return res, err
	}
	r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Response.Write(res.HTML)
	return
}
