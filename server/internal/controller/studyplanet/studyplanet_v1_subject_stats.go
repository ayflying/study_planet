package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// SubjectStats 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) SubjectStats(ctx context.Context, req *v1.SubjectStatsReq) (res *v1.SubjectStatsRes, err error) {
	return svc.StudyPlanet().SubjectStats(ctx, req)
}
