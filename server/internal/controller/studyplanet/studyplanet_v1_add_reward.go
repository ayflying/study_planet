package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// AddReward 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) AddReward(ctx context.Context, req *v1.AddRewardReq) (res *v1.AddRewardRes, err error) {
	return svc.StudyPlanet().AddReward(ctx, req)
}
