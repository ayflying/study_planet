package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"studyplanet/internal/model"
)

// ---------- 多邻国式练习场次：星级 + 连击 + XP ----------

// 连击奖分规则（多邻国式）：连击每达到 3/5/8 阶梯，给一次额外分。
var comboBonus = map[int]int{3: 2, 5: 4, 8: 6, 10: 8}

// CreateSession 开启一关：POST /api/sessions {subject, level, total, student_id}
func (s *Store) CreateSession(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	var body struct {
		Subject string `json:"subject"`
		Level   int    `json:"level"`
		Total   int    `json:"total"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	switch body.Subject {
	case "words", "reading", "math":
	default:
		s.fail(r, 400, "subject 需为 words/reading/math")
		return
	}
	if body.Level <= 0 {
		body.Level = 1
	}
	if body.Total <= 0 || body.Total > 50 {
		body.Total = 5
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	res, err := s.DB.Exec(
		"INSERT INTO practice_sessions(child_id,subject,level,total,created_at) VALUES(?,?,?,?,?)",
		cid, body.Subject, body.Level, body.Total, now,
	)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	id64, _ := res.LastInsertId()
	var ps model.PracticeSession
	if err := s.DB.Get(&ps, "SELECT id,child_id,subject,level,total,correct,max_combo,bonus,stars,finished,created_at FROM practice_sessions WHERE id=?", id64); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, ps)
}

// answerOutcome 单题作答后的实时反馈。
type answerOutcome struct {
	Correct    bool   `json:"correct"`
	Combo      int    `json:"combo"`            // 当前连击
	BasePoints int    `json:"base_points"`      // 本题基础得分（含连击加成）
	ComboBonus int    `json:"combo_bonus"`      // 本次触发的连击阶梯奖分
	Answer     string `json:"answer,omitempty"` // 正确答案（答错时反馈）
}

// recordAnswer 场次内记一笔作答并处理积分与连击；session 必须属于该学生且未结束。
func (s *Store) recordAnswer(r *ghttp.Request, sessionID int, refID int, correct bool, basePoints int, answer string) *answerOutcome {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return nil
	}
	var ps model.PracticeSession
	if err := s.DB.Get(&ps, "SELECT id,child_id,subject,level,total,correct,max_combo,bonus,stars,finished FROM practice_sessions WHERE id=?", sessionID); err != nil {
		s.fail(r, 404, "未找到该练习")
		return nil
	}
	if ps.ChildID != cid || ps.Finished == 1 {
		s.fail(r, 400, "该练习不可作答")
		return nil
	}

	// 当前连击 = 最近一次答错之后的连续正确数
	streak := s.streakOf(sessionID)

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

	now := time.Now().Format("2006-01-02 15:04:05")
	if _, err := s.DB.Exec(
		"INSERT INTO session_answers(session_id,ref_id,correct,points,combo,answered_at) VALUES(?,?,?,?,?,?)",
		sessionID, refID, boolToInt(correct), points, combo, now,
	); err != nil {
		s.fail(r, 500, err.Error())
		return nil
	}
	if correct {
		if _, err := s.DB.Exec("UPDATE practice_sessions SET correct=correct+1 WHERE id=?", sessionID); err != nil {
			s.fail(r, 500, err.Error())
			return nil
		}
		if combo > ps.MaxCombo {
			if _, err := s.DB.Exec("UPDATE practice_sessions SET max_combo=? WHERE id=?", combo, sessionID); err != nil {
				s.fail(r, 500, err.Error())
				return nil
			}
		}
		if points > 0 {
			reason := fmt.Sprintf("%s 第%d题", subjectName(ps.Subject), countAnswers(s, sessionID)+1)
			if combo >= 3 {
				reason += fmt.Sprintf(" 连击x%d", combo)
			}
			s.award(cid, points, reason)
		}
	} else if _, err := s.DB.Exec("UPDATE practice_sessions SET correct=correct WHERE id=?", sessionID); err != nil {
		s.fail(r, 500, err.Error())
		return nil
	}

	out := &answerOutcome{Correct: correct, Combo: combo, BasePoints: points, ComboBonus: comboBonusGot, Answer: answer}
	s.ok(r, out)
	return out
}

// countAnswers 该场已作答题数（用于流水备注）。
func countAnswers(s *Store, sessionID int) int {
	var n int
	_ = s.DB.Get(&n, "SELECT COUNT(*) FROM session_answers WHERE session_id=?", sessionID)
	return n
}

// streakOf 最近一次答错之后到现在的连续正确数。
func (s *Store) streakOf(sessionID int) int {
	var rows []struct {
		Correct int `db:"correct"`
	}
	if err := s.DB.Select(&rows, "SELECT correct FROM session_answers WHERE session_id=? ORDER BY id DESC LIMIT 50", sessionID); err != nil {
		return 0
	}
	n := 0
	for _, row := range rows {
		if row.Correct == 1 {
			n++
		} else {
			break
		}
	}
	return n
}

// lastWrongOffset 占位避免未使用告警（保留函数签名以便扩展）。
func lastWrongOffset(sessionID int) int { return 50 }

// subjectName 科目中文名（积分流水用）。
func subjectName(subject string) string {
	switch subject {
	case "words":
		return "单词"
	case "reading":
		return "阅读"
	case "math":
		return "数学"
	}
	return "练习"
}

// FinishSession 结算一关：POST /api/sessions/:id/finish
// 星级规则：正确率 >=90% 三星、>=70% 两星、>=50% 一星、否则零星；
// 三星额外奖励 10 分，两星 5 分。同一关只结算一次。
func (s *Store) FinishSession(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	id := s.idParam(r)
	var ps model.PracticeSession
	if err := s.DB.Get(&ps, "SELECT id,child_id,subject,level,total,correct,max_combo,bonus,stars,finished FROM practice_sessions WHERE id=?", id); err != nil {
		s.fail(r, 404, "未找到该练习")
		return
	}
	if ps.ChildID != cid {
		s.fail(r, 400, "无权结算他人练习")
		return
	}
	if ps.Finished == 1 {
		s.ok(r, map[string]interface{}{"already": true, "stars": ps.Stars, "bonus": ps.Bonus, "max_combo": ps.MaxCombo})
		return
	}
	ratio := 0.0
	if ps.Total > 0 {
		ratio = float64(ps.Correct) / float64(ps.Total)
	}
	stars := 0
	bonus := 0
	switch {
	case ratio >= 0.9:
		stars, bonus = 3, 10
	case ratio >= 0.7:
		stars, bonus = 2, 5
	case ratio >= 0.5:
		stars = 1
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	if _, err := s.DB.Exec(
		"UPDATE practice_sessions SET stars=?, bonus=?, finished=1, finished_at=? WHERE id=?",
		stars, bonus, now, id,
	); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	if bonus > 0 {
		label := "关卡完成"
		if stars == 3 {
			label = "三星通关"
		}
		s.award(cid, bonus, fmt.Sprintf("%s 奖励:+%d", label, bonus))
	}
	s.ok(r, map[string]interface{}{
		"stars": stars, "bonus": bonus, "max_combo": ps.MaxCombo,
		"correct": ps.Correct, "total": ps.Total,
	})
}

// ListSessions 学生最近的练习记录（家长端统计用）：GET /api/sessions?student_id=
func (s *Store) ListSessions(r *ghttp.Request) {
	cid := s.resolveChild(r)
	if cid < 0 {
		s.fail(r, 404, "学生不存在")
		return
	}
	var ss []model.PracticeSession
	q := "SELECT id,child_id,subject,level,total,correct,max_combo,bonus,stars,finished,created_at,COALESCE(finished_at,'') AS finished_at FROM practice_sessions WHERE child_id=?"
	args := []interface{}{cid}
	if lv := r.GetQuery("level").String(); lv != "" {
		q += " AND level=?"
		args = append(args, lv)
	}
	if sub := r.GetQuery("subject").String(); sub != "" {
		q += " AND subject=?"
		args = append(args, sub)
	}
	q += " ORDER BY id DESC LIMIT 50"
	if err := s.DB.Select(&ss, q, args...); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, ss)
}

// 答题接口改造辅助：把原「直接判分」的三个接口切换为「连击+场次」模式。
// StartOrReuseSession 保证进行中的场次存在（前端进入关卡时调用）。
func levelParam(r *ghttp.Request, def int) int {
	if v := r.GetQuery("level").Int(); v > 0 {
		return v
	}
	return def
}

func itoa(v int) string { return strconv.Itoa(v) }

var _ = itoa // 预留

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
