package v1

import "github.com/gogf/gf/v2/frame/g"

// AuthModeReq 登录模式查询：casdoor 或 pin。
type AuthModeReq struct {
	g.Meta `path:"/parent/auth-mode" method:"get" tags:"ParentAuth" summary:"查询家长登录模式"`
}
type AuthModeRes struct {
	Mode        string `json:"mode"`
	LoginURL    string `json:"login_url,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
	Application string `json:"application,omitempty"`
}

// ParentLoginReq 家长 PIN 登录（Casdoor 启用时禁用）。
type ParentLoginReq struct {
	g.Meta `path:"/parent/login" method:"post" tags:"ParentAuth" summary:"家长 PIN 登录"`
	Pin    string `json:"pin" v:"required"`
}
type ParentLoginRes struct {
	Token string `json:"token"`
}

// CasdoorLoginReq 跳转 Casdoor 授权页（302）。
type CasdoorLoginReq struct {
	g.Meta      `path:"/parent/casdoor/login" method:"get,all" tags:"ParentAuth" summary:"Casdoor 登录跳转"`
	// RedirectURI 控制器按请求推导的回调地址。
	RedirectURI string `json:"-"`
}
type CasdoorLoginRes struct {
	// Location 为 Casdoor 授权页完整地址，由 logic 计算返回，
	// 控制器据此执行 302 跳转（HTTP 语义保留在边缘层）。
	Location string `json:"-"`
}

// CasdoorCallbackReq Casdoor 授权码回调：写 token 到 localStorage 后回首页。
type CasdoorCallbackReq struct {
	g.Meta     `path:"/parent/casdoor/callback" method:"get,all" tags:"ParentAuth" summary:"Casdoor 回调"`
	Code       string `json:"code"`
	// RedirectURI 控制器按请求推导的回调地址（logic 用来向 Casdoor 换 token）。
	RedirectURI string `json:"-"`
	// ClientIP / TLS 语义同上，保留 HTTP 边缘信息由控制器注入。
	IsTLS bool `json:"-"`
}
type CasdoorCallbackRes struct {
	// HTML 为登录完成后的回跳页（写入 localStorage 后跳首页），
	// 控制器以 text/html 输出（HTTP 语义保留在边缘层）。
	HTML string `json:"-"`
}

// SetPinReq 修改家长 PIN（家长鉴权）。
type SetPinReq struct {
	g.Meta `path:"/parent/set-pin" method:"post" tags:"ParentAuth" summary:"修改家长 PIN"`
	Pin    string `json:"pin" v:"required|length:4,32"`
}
type SetPinRes struct {
	OK bool `json:"ok"`
}
