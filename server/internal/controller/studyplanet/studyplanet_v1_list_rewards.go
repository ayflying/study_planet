package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// ListRewards 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) ListRewards(ctx context.Context, req *v1.ListRewardsReq) (res *v1.ListRewardsRes, err error) {
	return svc.StudyPlanet().ListRewards(ctx, req)
}
