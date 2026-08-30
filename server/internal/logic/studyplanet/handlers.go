// Package studyplanet 逻辑层主文件：健康检查与废弃接口占位。
// 业务实现按模块拆分：words/reading/math/tasks/points/students/practice/pet/content 等。
package studyplanet

import (
	"context"
	"time"

	v1 "studyplanet/api/studyplanet/v1"
)

// Health 健康检查：返回运行状态与版本号。
func (s *sStudyPlanet) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	return &v1.HealthRes{
		Status:  "ok",
		Time:    time.Now().Format(time.RFC3339),
		App:     "studyplanet",
		Version: CurrentVersion(),
	}, nil
}

// ParentLogin 家长 PIN 登录（已废弃：多家长架构下 PIN 无法区分家长身份，强制走 Casdoor SSO）。
func (s *sStudyPlanet) ParentLogin(ctx context.Context, req *v1.ParentLoginReq) (res *v1.ParentLoginRes, err error) {
	return nil, errAuth("PIN 登录已停用，请使用 Casdoor SSO 登录")
}

// SetPin 修改家长 PIN（已废弃：PIN 登录停用，保留接口返回提示避免旧前端报 404）。
func (s *sStudyPlanet) SetPin(ctx context.Context, req *v1.SetPinReq) (res *v1.SetPinRes, err error) {
	return nil, errParam("PIN 登录已停用，无需设置 PIN")
}
