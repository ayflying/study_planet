package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// WordProgress 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) WordProgress(ctx context.Context, req *v1.WordProgressReq) (res *v1.WordProgressRes, err error) {
	return svc.StudyPlanet().WordProgress(ctx, req)
}
