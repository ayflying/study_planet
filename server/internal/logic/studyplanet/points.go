// Package studyplanet 积分与奖励模块：积分汇总/流水、奖励列表、兑换与家长确认。
package studyplanet

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

// pointsTotal 指定学生当前积分总量（points_log 汇总）。
func (s *sStudyPlanet) pointsTotal(ctx context.Context, childID int) (int, error) {
	v, err := daoPointsLog.Ctx(ctx).Fields("COALESCE(SUM(delta),0) AS total").Where("child_id", childID).Value()
	if err != nil {
		return 0, gerror.Wrap(err, "查询积分失败")
	}
	return v.Int(), nil
}

// PointsSummary 积分汇总（公开接口）。
func (s *sStudyPlanet) PointsSummary(ctx context.Context, req *v1.PointsSummaryReq) (res *v1.PointsSummaryRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	total, err := s.pointsTotal(ctx, cid)
	if err != nil {
		return nil, err
	}
	v, err := daoPointsLog.Ctx(ctx).
		Fields("COALESCE(SUM(delta),0) AS total").
		Where("child_id", cid).
		Where("DATE(created_at)=?", todayStr()).Value()
	if err != nil {
		return nil, gerror.Wrap(err, "查询今日积分失败")
	}
	return &v1.PointsSummaryRes{Total: total, TodayEarned: v.Int(), StudentID: cid}, nil
}

// PointsLog 积分流水（最近 100 条，公开接口）。
func (s *sStudyPlanet) PointsLog(ctx context.Context, req *v1.PointsLogReq) (res *v1.PointsLogRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	all, err := daoPointsLog.Ctx(ctx).Where("child_id", cid).OrderDesc("id").Limit(100).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询积分流水失败")
	}
	out := make(v1.PointsLogRes, 0, len(all))
	for _, r := range all {
		out = append(out, v1.PointsLogItem{
			ID:        r["id"].Int(),
			ChildID:   r["child_id"].Int(),
			Delta:     r["delta"].Int(),
			Reason:    r["reason"].String(),
			CreatedAt: r["created_at"].String(),
		})
	}
	return &out, nil
}

// ListRewards 奖励列表。
// 学生端（无家长 token）：只展示当前孩子所属家长的奖励，家长之间互相隔离；
// 未指定 student_id 时返回空列表（不泄露任何家长数据）。
func (s *sStudyPlanet) ListRewards(ctx context.Context, req *v1.ListRewardsReq) (res *v1.ListRewardsRes, err error) {
	out := make(v1.ListRewardsRes, 0, 8)
	if req.StudentID > 0 {
		ownerID := s.childParentID(ctx, req.StudentID)
		if ownerID > 0 {
			all, err := daoRewards.Ctx(ctx).
				Where("parent_id", ownerID).
				Where("status", "active").
				Order("cost_points").All()
			if err != nil {
				return nil, gerror.Wrap(err, "查询奖励失败")
			}
			for _, r := range all {
				out = append(out, v1.Reward{
					ID:         r["id"].Int(),
					Name:       r["name"].String(),
					CostPoints: r["cost_points"].Int(),
					Status:     r["status"].String(),
				})
			}
		}
	}
	return &out, nil
}

// Redeem 学生兑换奖励（公开接口，需家长确认）。
// 只能兑换当前孩子所属家长上架的奖励，跨家庭兑换直接拒绝。
func (s *sStudyPlanet) Redeem(ctx context.Context, req *v1.RedeemReq) (res *v1.RedeemRes, err error) {
	rw, err := daoRewards.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询奖励失败")
	}
	if rw.IsEmpty() {
		return nil, errNotFound("未找到该奖励")
	}
	if rw["status"].String() != "active" {
		return nil, errParam("该奖励暂不可用")
	}
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	// 归属校验：奖励必须属于该孩子所在家庭
	if ownerID := s.childParentID(ctx, cid); ownerID == 0 || ownerID != rw["parent_id"].Int() {
		return nil, errForbidden("该奖励不属于当前家庭")
	}
	total, err := s.pointsTotal(ctx, cid)
	if err != nil {
		return nil, err
	}
	if total < rw["cost_points"].Int() {
		return nil, errParam("积分不足")
	}
	if _, err := daoRedempt.Ctx(ctx).Data(doRedempt{
		RewardId:    req.ID,
		ChildId:     cid,
		Status:      "pending",
		RequestedAt: gtime.Now(),
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "提交兑换失败")
	}
	return &v1.RedeemRes{OK: true, Pending: true, Message: "已提交兑换，等待家长确认"}, nil
}

// AddReward 添加奖励（家长鉴权）：奖励归属当前家长，仅自家孩子可见可兑。
func (s *sStudyPlanet) AddReward(ctx context.Context, req *v1.AddRewardReq) (res *v1.AddRewardRes, err error) {
	if err := s.ensureParentAuth(ctx); err != nil {
		return nil, err
	}
	if err := ensureName(req.Name, "请填写奖励名称"); err != nil {
		return nil, err
	}
	if _, err := daoRewards.Ctx(ctx).Data(doRewards{
		Name:       req.Name,
		CostPoints: req.CostPoints,
		Status:     "active",
		ParentId:   ctxParentID(ctx),
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "添加奖励失败")
	}
	return &v1.AddRewardRes{OK: true}, nil
}

// ConfirmRedemption 家长确认兑换（家长鉴权）：兑换必须属于当前家长的学生。
func (s *sStudyPlanet) ConfirmRedemption(ctx context.Context, req *v1.ConfirmRedemptionReq) (res *v1.ConfirmRedemptionRes, err error) {
	rd, err := daoRedempt.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询兑换失败")
	}
	if rd.IsEmpty() {
		return nil, errNotFound("未找到该兑换")
	}
	if err := s.ensureChildOwned(ctx, rd["child_id"].Int()); err != nil {
		return nil, err
	}
	if rd["status"].String() != "pending" {
		return nil, errAuth("该兑换不在待确认状态")
	}
	rw, err := daoRewards.Ctx(ctx).Where("id", rd["reward_id"]).One()
	if err != nil || rw.IsEmpty() {
		return nil, errNotFound("未找到对应奖励")
	}
	if _, err := daoRedempt.Ctx(ctx).Where("id", req.ID).Data(doRedempt{
		Status:      "confirmed",
		ConfirmedAt: gtime.Now(),
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "确认兑换失败")
	}
	if _, err := daoRewards.Ctx(ctx).Where("id", rd["reward_id"]).Data(doRewards{
		Status: "redeemed",
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "更新奖励状态失败")
	}
	s.award(rd["child_id"].Int(), -rw["cost_points"].Int(), "兑换:"+rw["name"].String())
	return &v1.ConfirmRedemptionRes{OK: true}, nil
}
