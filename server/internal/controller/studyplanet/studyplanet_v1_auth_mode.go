package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// AuthMode 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) AuthMode(ctx context.Context, req *v1.AuthModeReq) (res *v1.AuthModeRes, err error) {
	return svc.StudyPlanet().AuthMode(ctx, req)
}
