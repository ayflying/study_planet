package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// PointsLog 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) PointsLog(ctx context.Context, req *v1.PointsLogReq) (res *v1.PointsLogRes, err error) {
	return svc.StudyPlanet().PointsLog(ctx, req)
}
