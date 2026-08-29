package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// FinishSession 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) FinishSession(ctx context.Context, req *v1.FinishSessionReq) (res *v1.FinishSessionRes, err error) {
	return svc.StudyPlanet().FinishSession(ctx, req)
}
