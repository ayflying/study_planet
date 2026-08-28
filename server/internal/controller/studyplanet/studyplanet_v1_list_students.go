package studyplanet

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ListStudents 控制器转发：logic 实现直接写响应（保持历史 JSON 结构），
// 返回 CodeOK 抑制 GF 默认响应包装，保证对外 JSON 与历史完全一致。
func (c *ControllerV1) ListStudents(ctx context.Context, req *v1.ListStudentsReq) (res *v1.ListStudentsRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidRequest)
	}
	svc.StudyPlanet().ListStudents(r)
	return nil, gerror.NewCode(gcode.CodeOK)
}
