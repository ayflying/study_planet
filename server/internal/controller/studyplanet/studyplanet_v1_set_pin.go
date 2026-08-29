package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// SetPin 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) SetPin(ctx context.Context, req *v1.SetPinReq) (res *v1.SetPinRes, err error) {
	return svc.StudyPlanet().SetPin(ctx, req)
}
