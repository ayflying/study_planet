package studyplanet

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

// ---------- 多邻国式练习场次：星级 + 连击 + XP ----------

// 连击奖分规则（多邻国式）：连击每达到 3/5/8/10 阶梯，给一次额外分。
var comboBonus = map[int]int{3: 2, 5: 4, 8: 6, 10: 8}

// answerOutcome 单题作答后的实时反馈（场次模式）。
type answerOutcome struct {
	Correct    bool   // 本题是否正确
	Combo      int    // 当前连击
	BasePoints int    // 本题基础得分（含连击加成）
	ComboBonus int    // 本次触发的连击阶梯奖分
	Answer     string // 正确答案（答错时反馈）
	Review     bool   // 本题是否为错题巩固题
	XP         int    // 本次作答获得的经验值
}

// fillOutcome 把作答反馈写入接口响应字段（correct 可为 nil）。
func fillOutcome(combo, basePoints, comboBonusGot, review *int, xp *int, correct *bool, out *answerOutcome) {
	if out == nil {
		return
	}
	if combo != nil {
		*combo = out.Combo
	}
	if basePoints != nil {
		*basePoints = out.BasePoints
	}
	if comboBonusGot != nil {
		*comboBonusGot = out.ComboBonus
	}
	if review != nil {
		if out.Review {
			*review = 1
		} else {
			*review = 0
		}
	}
	if xp != nil {
		*xp = out.XP
	}
	if correct != nil {
		*correct = out.Correct
	}
}

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

// recordAnswer 场次内记一笔作答并处理积分与连击；session 必须属于该学生且未结束。
// reviewRefs 标记哪些 ref_id 来自错题本（答对则消除错题）；可为 nil。
func (s *sStudyPlanet) recordAnswer(ctx context.Context, childID, sessionID, refID int, correct bool, basePoints int, answer string, reviewRefs map[int]bool) (*answerOutcome, error) {
	if childID < 0 {
		return nil, errNotFound("学生不存在")
	}
	ps, err := daoSessions.Ctx(ctx).Where("id", sessionID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询练习场次失败")
	}
	if ps.IsEmpty() {
		return nil, errNotFound("未找到该练习")
	}
	if ps["child_id"].Int() != childID || ps["finished"].Int() == 1 {
		return nil, errAuth("该练习不可作答")
	}

	// 当前连击 = 最近一次答错之后的连续正确数
	streak, err := s.streakOf(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	combo := streak
	comboBonusGot := 0
	points := 0
	if correct {
		combo = streak + 1
		points = basePoints
		if b, ok := comboBonus[combo]; ok {
			comboBonusGot = b
			points += b
		}
	}

	review := isReviewRef(reviewRefs, refID)
	if _, err := daoAnswers.Ctx(ctx).Data(doAnswers{
		SessionId:  sessionID,
		RefId:      refID,
		Correct:    boolToInt(correct),
		Points:     points,
		Combo:      combo,
		AnsweredAt: gtime.Now(),
	}).Insert(); err != nil {
		return nil, gerror.Wrap(err, "记录作答失败")
	}
	if correct {
		if _, err := daoSessions.Ctx(ctx).Where("id", sessionID).Increment("correct", 1); err != nil {
			return nil, gerror.Wrap(err, "更新正确数失败")
		}
		if combo > ps["max_combo"].Int() {
			if _, err := daoSessions.Ctx(ctx).Where("id", sessionID).Data(doSessions{MaxCombo: combo}).Update(); err != nil {
				return nil, gerror.Wrap(err, "更新连击失败")
			}
		}
		xp := 0
		if points > 0 {
			reason := fmt.Sprintf("%s 第%d题", subjectName(ps["subject"].String()), s.countAnswers(ctx, sessionID)+1)
			if combo >= 3 {
				reason += fmt.Sprintf(" 连击x%d", combo)
			}
			s.award(childID, points, reason)
			// 正反馈：答题得分同步转化为经验值（1:1）
			xp = points
			s.addXP(ctx, childID, xp)
		}
		// 错题巩固答对 → 从错题本消除
		if review {
			s.resolveWrong(ctx, childID, ps["subject"].String(), refID)
		}
		return &answerOutcome{Correct: correct, Combo: combo, BasePoints: points, ComboBonus: comboBonusGot, Answer: answer, Review: review, XP: xp}, nil
	}
	// 答错 → 登记错题本
	s.recordWrong(ctx, childID, ps["subject"].String(), refID)
	return &answerOutcome{Correct: correct, Combo: combo, BasePoints: points, ComboBonus: comboBonusGot, Answer: answer, Review: review, XP: 0}, nil
}

// countAnswers 该场已作答题数（用于积分流水备注）。
func (s *sStudyPlanet) countAnswers(ctx context.Context, sessionID int) int {
	n, err := daoAnswers.Ctx(ctx).Where("session_id", sessionID).Count()
	if err != nil {
		return 0
	}
	return n
}

// streakOf 最近一次答错之后到现在的连续正确数。
func (s *sStudyPlanet) streakOf(ctx context.Context, sessionID int) (int, error) {
	all, err := daoAnswers.Ctx(ctx).Fields("correct").Where("session_id", sessionID).OrderDesc("id").Limit(50).All()
	if err != nil {
		return 0, gerror.Wrap(err, "查询作答记录失败")
	}
	n := 0
	for _, r := range all {
		if r["correct"].Int() == 1 {
			n++
			continue
		}
		break
	}
	return n, nil
}

// subjectName 科目中文名（积分流水用）。
func subjectName(subject string) string {
	switch strings.TrimSpace(subject) {
	case "words", "english":
		return "英语"
	case "reading", "chinese":
		return "语文"
	case "math":
		return "数学"
	case "physics":
		return "物理"
	case "chemistry":
		return "化学"
	case "biology":
		return "生物"
	case "history":
		return "历史"
	case "geography":
		return "地理"
	}
	return "练习"
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

// boolToInt 布尔 → 数据库 0/1。
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// 保持 g 包引用（扩展查询用），避免未使用导入。
var _ = g.DB
