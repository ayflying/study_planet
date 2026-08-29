package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// PointsSummary 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) PointsSummary(ctx context.Context, req *v1.PointsSummaryReq) (res *v1.PointsSummaryRes, err error) {
	return svc.StudyPlanet().PointsSummary(ctx, req)
}
