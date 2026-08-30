// Package studyplanet 场次生命周期：开启一关、结算星级、查询记录。
package studyplanet

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

// CreateSession 开启一关（多邻国式关卡）。
func (s *sStudyPlanet) CreateSession(ctx context.Context, req *v1.CreateSessionReq) (res *v1.CreateSessionRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	// 动态内容库后学科是开放集合（english/chinese/math/physics...），不能再限 words/reading/math。
	// 校验：非空 + 长度合法 + 必须是内容库中已启用学科（兼容旧值 words/reading）。
	if req.Subject == "" || len(req.Subject) > 32 {
		return nil, errParam("subject 不能为空")
	}
	if req.Subject != "words" && req.Subject != "reading" {
		cnt, err := daoSubjects.Ctx(ctx).Where("code", req.Subject).Where("enabled", 1).Count()
		if err != nil {
			return nil, gerror.Wrap(err, "校验学科失败")
		}
		if cnt == 0 {
			return nil, errParam("subject 需为已启用的学科 code（如 english/chinese/math）")
		}
	}
	level, total := req.Level, req.Total
	if level <= 0 {
		level = 1
	}
	if total <= 0 || total > 50 {
		total = 5
	}
	id, err := daoSessions.Ctx(ctx).Data(doSessions{
		ChildId:   cid,
		Subject:   req.Subject,
		Level:     level,
		Total:     total,
		CreatedAt: gtime.Now(),
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.Wrap(err, "创建练习场次失败")
	}
	rec, err := daoSessions.Ctx(ctx).Where("id", id).One()
	if err != nil || rec.IsEmpty() {
		return nil, gerror.New("查询练习场次失败")
	}
	ps := v1.CreateSessionRes(sessionOf(rec))
	return &ps, nil
}

// FinishSession 结算一关。
// 星级规则：正确率 >=90% 三星、>=70% 两星、>=50% 一星、否则零星；
// 三星额外奖励 10 分，两星 5 分。同一关只结算一次。
func (s *sStudyPlanet) FinishSession(ctx context.Context, req *v1.FinishSessionReq) (res *v1.FinishSessionRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	ps, err := daoSessions.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询练习场次失败")
	}
	if ps.IsEmpty() {
		return nil, errNotFound("未找到该练习")
	}
	if ps["child_id"].Int() != cid {
		return nil, errAuth("无权结算他人练习")
	}
	if ps["finished"].Int() == 1 {
		return &v1.FinishSessionRes{
			Stars:    ps["stars"].Int(),
			Bonus:    ps["bonus"].Int(),
			MaxCombo: ps["max_combo"].Int(),
			Already:  1,
		}, nil
	}
	total, correct := ps["total"].Int(), ps["correct"].Int()
	ratio := 0.0
	if total > 0 {
		ratio = float64(correct) / float64(total)
	}
	stars, bonus := 0, 0
	switch {
	case ratio >= 0.9:
		stars, bonus = 3, 10
	case ratio >= 0.7:
		stars, bonus = 2, 5
	case ratio >= 0.5:
		stars = 1
	}
	if _, err := daoSessions.Ctx(ctx).Where("id", req.ID).Data(doSessions{
		Stars:      stars,
		Bonus:      bonus,
		Finished:   1,
		FinishedAt: gtime.Now(),
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "结算练习失败")
	}
	// 星级奖励：奖分（积分）+ 经验值（1 星 10 / 2 星 20 / 3 星 35）
	xpGain := map[int]int{1: 10, 2: 20, 3: 35}[stars]
	if xpGain > 0 {
		s.addXP(ctx, cid, xpGain)
	}
	if bonus > 0 {
		label := "关卡完成"
		if stars == 3 {
			label = "三星通关"
		}
		s.award(cid, bonus, fmt.Sprintf("%s 奖励:+%d", label, bonus))
	}
	// 学习探索地图完成一关可能掉落零食
	s.addSnackDrop(ctx, cid)
	return &v1.FinishSessionRes{
		Stars:    stars,
		Bonus:    bonus,
		MaxCombo: ps["max_combo"].Int(),
		Correct:  correct,
		Total:    total,
		XPGained: xpGain,
	}, nil
}

// ListSessions 学生最近的练习记录（家长端统计用）。
func (s *sStudyPlanet) ListSessions(ctx context.Context, req *v1.ListSessionsReq) (res *v1.ListSessionsRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	m := daoSessions.Ctx(ctx).Where("child_id", cid)
	if req.Level != "" {
		m = m.Where("level", req.Level)
	}
	if req.Subject != "" {
		m = m.Where("subject", req.Subject)
	}
	all, err := m.OrderDesc("id").Limit(50).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询练习记录失败")
	}
	out := make(v1.ListSessionsRes, 0, len(all))
	for _, r := range all {
		out = append(out, sessionOf(r))
	}
	return &out, nil
}
