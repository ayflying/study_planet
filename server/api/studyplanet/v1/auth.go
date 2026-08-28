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
	g.Meta `path:"/parent/casdoor/login" method:"get,all" tags:"ParentAuth" summary:"Casdoor 登录跳转"`
}
type CasdoorLoginRes struct{}

// CasdoorCallbackReq Casdoor 授权码回调：写 token 到 localStorage 后回首页。
type CasdoorCallbackReq struct {
	g.Meta `path:"/parent/casdoor/callback" method:"get,all" tags:"ParentAuth" summary:"Casdoor 回调"`
	Code   string `json:"code"`
}
type CasdoorCallbackRes struct{}

// SetPinReq 修改家长 PIN（家长鉴权）。
type SetPinReq struct {
	g.Meta `path:"/parent/set-pin" method:"post" tags:"ParentAuth" summary:"修改家长 PIN"`
	Pin    string `json:"pin" v:"required|length:4,32"`
}
type SetPinRes struct {
	OK bool `json:"ok"`
}
