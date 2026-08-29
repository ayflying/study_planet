package studyplanet

import (
	"context"

	v1 "studyplanet/api/studyplanet/v1"
	svc "studyplanet/internal/service"
)

// PetGet 控制器：我的宠物（首次访问自动领取）。
func (c *ControllerV1) PetGet(ctx context.Context, req *v1.PetGetReq) (res *v1.PetGetRes, err error) {
	return svc.StudyPlanet().PetGet(ctx, req)
}

// PetFeed 控制器：投喂宠物。
func (c *ControllerV1) PetFeed(ctx context.Context, req *v1.PetFeedReq) (res *v1.PetFeedRes, err error) {
	return svc.StudyPlanet().PetFeed(ctx, req)
}

// PetRename 控制器：宠物改名。
func (c *ControllerV1) PetRename(ctx context.Context, req *v1.PetRenameReq) (res *v1.PetRenameRes, err error) {
	return svc.StudyPlanet().PetRename(ctx, req)
}

// PetTick 控制器：宠物状态惰性结算。
func (c *ControllerV1) PetTick(ctx context.Context, req *v1.PetDecayReq) (res *v1.PetDecayRes, err error) {
	return svc.StudyPlanet().PetTick(ctx, req)
}

// PetFoods 控制器：食物列表。
func (c *ControllerV1) PetFoods(ctx context.Context, req *v1.PetFoodsReq) (res *v1.PetFoodsRes, err error) {
	return svc.StudyPlanet().PetFoods(ctx, req)
}

// BattleRank 控制器：对战段位榜。
func (c *ControllerV1) BattleRank(ctx context.Context, req *v1.BattleRankReq) (res *v1.BattleRankRes, err error) {
	return svc.StudyPlanet().BattleRank(ctx, req)
}

// BattleHistory 控制器：我的历史对战。
func (c *ControllerV1) BattleHistory(ctx context.Context, req *v1.BattleHistoryReq) (res *v1.BattleHistoryRes, err error) {
	return svc.StudyPlanet().BattleHistory(ctx, req)
}
