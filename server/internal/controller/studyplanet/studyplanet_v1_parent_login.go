package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ParentLogin 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ParentLogin(ctx context.Context, req *v1.ParentLoginReq) (res *v1.ParentLoginRes, err error) {
	return svc.StudyPlanet().ParentLogin(ctx, req)
}
