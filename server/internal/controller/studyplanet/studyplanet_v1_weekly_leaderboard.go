package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// WeeklyLeaderboard 控制器：仅做参数绑定后的转发，业务逻辑一律在 logic 层。
func (c *ControllerV1) WeeklyLeaderboard(ctx context.Context, req *v1.WeeklyLeaderboardReq) (res *v1.WeeklyLeaderboardRes, err error) {
	return svc.StudyPlanet().WeeklyLeaderboard(ctx, req)
}
