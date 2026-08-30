// Package studyplanet 场次内作答：连击计分、XP 转化、错题登记。
package studyplanet

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

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

// boolToInt 布尔 → 数据库 0/1。
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
