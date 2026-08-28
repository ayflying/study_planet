// Package v1 StudyPlanet API v1 接口定义（gf gen ctrl 的输入源）。
// 命名规范：操作+Req / 操作+Res；g.Meta 携带路由路径与方法。
// 注意：路径不含 /api 前缀（由 router 的 /api 分组统一挂载），与历史 URL 保持一致。
package v1

import "github.com/gogf/gf/v2/frame/g"

// HealthReq 健康检查。
type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"System" summary:"健康检查"`
}
type HealthRes struct {
	Status  string `json:"status"`
	Time    string `json:"time"`
	App     string `json:"app"`
	Version string `json:"version"`
}
