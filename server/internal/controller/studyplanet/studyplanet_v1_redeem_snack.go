package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// RedeemSnack 控制器：星星兑换零食。
func (c *ControllerV1) RedeemSnack(ctx context.Context, req *v1.RedeemSnackReq) (res *v1.RedeemSnackRes, err error) {
	return svc.StudyPlanet().RedeemSnack(ctx, req)
}