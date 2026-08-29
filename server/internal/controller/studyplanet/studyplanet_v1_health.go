package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// Health 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	return svc.StudyPlanet().Health(ctx, req)
}
