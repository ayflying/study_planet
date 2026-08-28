package studyplanet

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// Health 控制器转发：logic 实现直接写响应（保持历史 JSON 结构），
// 返回 CodeOK 抑制 GF 默认响应包装，保证对外 JSON 与历史完全一致。
func (c *ControllerV1) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidRequest)
	}
	svc.StudyPlanet().Health(r)
	return nil, gerror.NewCode(gcode.CodeOK)
}
