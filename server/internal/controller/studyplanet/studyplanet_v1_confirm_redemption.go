package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ConfirmRedemption 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ConfirmRedemption(ctx context.Context, req *v1.ConfirmRedemptionReq) (res *v1.ConfirmRedemptionRes, err error) {
	return svc.StudyPlanet().ConfirmRedemption(ctx, req)
}
