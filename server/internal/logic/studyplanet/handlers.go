package studyplanet

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

// ---------- 通用读取助手（Record → api 结构） ----------

// studentOf children 记录 → 学生档案。
func studentOf(r gdb.Record) v1.Student {
	return v1.Student{
		ID:        r["id"].Int(),
		Name:      r["name"].String(),
		Username:  r["username"].String(),
		Avatar:    r["avatar"].String(),
		Grade:     r["grade"].Int(),
		CreatedAt: r["created_at"].String(),
	}
}

// sessionOf practice_sessions 记录 → 练习场次（时间列可空统一转字符串）。
func sessionOf(r gdb.Record) v1.PracticeSession {
	return v1.PracticeSession{
		ID:         r["id"].Int(),
		ChildID:    r["child_id"].Int(),
		Subject:    r["subject"].String(),
		Level:      r["level"].Int(),
		Total:      r["total"].Int(),
		Correct:    r["correct"].Int(),
		MaxCombo:   r["max_combo"].Int(),
		Bonus:      r["bonus"].Int(),
		Stars:      r["stars"].Int(),
		Finished:   r["finished"].Int(),
		CreatedAt:  r["created_at"].String(),
		FinishedAt: r["finished_at"].String(),
	}
}

// taskOf tasks 记录 → 任务（due_date/completed_at 可空，空值输出 ""）。
func taskOf(r gdb.Record) v1.Task {
	due, completed := "", ""
	if !r["due_date"].IsNil() {
		due = r["due_date"].GTime().Format("Y-m-d")
	}
	if !r["completed_at"].IsNil() {
		completed = r["completed_at"].GTime().Format("Y-m-d H:i:s")
	}
	return v1.Task{
		ID:          r["id"].Int(),
		Title:       r["title"].String(),
		Type:        r["type"].String(),
		DueDate:     due,
		Points:      r["points"].Int(),
		Status:      r["status"].String(),
		CreatedAt:   r["created_at"].String(),
		CompletedAt: completed,
	}
}

// resolveChild 解析本次请求对应的学生 id：来自 req 的 student_id，缺省为 1（兼容旧客户端）。
// 学生不存在时返回 -1。
// requireOwner=true（家长端写操作）时校验学生归属：parentID 与孩子 parent_id 不一致即拒绝，
// NULL 归属（历史数据尚未被接管）同样拒绝，防止越权操作他人孩子。
func (s *sStudyPlanet) resolveChild(ctx context.Context, studentID int) (int, error) {
	id := studentID
	if id <= 0 {
		id = 1
	}
	rec, err := daoChildren.Ctx(ctx).Fields("id,parent_id").Where("id", id).One()
	if err != nil {
		return -1, gerror.Wrap(err, "查询学生失败")
	}
	if rec.IsEmpty() {
		return -1, nil
	}
	return id, nil
}

// childParentID 查询孩子的归属家长 id（孩子不存在返回 0）。
func (s *sStudyPlanet) childParentID(ctx context.Context, childID int) int {
	v, err := daoChildren.Ctx(ctx).Fields("parent_id").Where("id", childID).Value()
	if err != nil || v == nil || v.IsNil() {
		return 0
	}
	return v.Int()
}

// ensureChildOwned 家长端操作前校验孩子归属：不属于自己的孩子一律拒绝。
func (s *sStudyPlanet) ensureChildOwned(ctx context.Context, childID int) error {
	parentID := ctxParentID(ctx)
	if parentID <= 0 {
		return errAuth("登录状态缺少家长身份，请重新登录")
	}
	owner := s.childParentID(ctx, childID)
	if owner == 0 {
		return errNotFound("学生不存在")
	}
	if owner != parentID {
		return errForbidden("无权操作该学生")
	}
	return nil
}

// pointsTotal 指定学生当前积分总量（points_log 汇总）。
func (s *sStudyPlanet) pointsTotal(ctx context.Context, childID int) (int, error) {
	v, err := daoPointsLog.Ctx(ctx).Fields("COALESCE(SUM(delta),0) AS total").Where("child_id", childID).Value()
	if err != nil {
		return 0, gerror.Wrap(err, "查询积分失败")
	}
	return v.Int(), nil
}

// ---------- 健康检查 ----------

// Health 健康检查：返回运行状态与版本号。
func (s *sStudyPlanet) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	return &v1.HealthRes{
		Status:  "ok",
		Time:    time.Now().Format(time.RFC3339),
		App:     "studyplanet",
		Version: CurrentVersion(),
	}, nil
}

// ---------- 单词卡片 ----------

// ListWords 单词列表（可选按 level 过滤）。
func (s *sStudyPlanet) ListWords(ctx context.Context, req *v1.ListWordsReq) (res *v1.ListWordsRes, err error) {
	m := daoWords.Ctx(ctx).Order("level", "id")
	if req.Level != "" {
		m = m.Where("level", req.Level)
	}
	all, err := m.All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询单词失败")
	}
	out := make(v1.ListWordsRes, 0, len(all))
	for _, r := range all {
		out = append(out, v1.Word{
			ID:        r["id"].Int(),
			Level:     r["level"].Int(),
			Word:      r["word"].String(),
			Meaning:   r["meaning"].String(),
			Phonetic:  r["phonetic"].String(),
			Example:   r["example"].String(),
			CreatedAt: r["created_at"].String(),
		})
	}
	return &out, nil
}

// WordDetail 单词详情 + 当前学生掌握状态。
func (s *sStudyPlanet) WordDetail(ctx context.Context, req *v1.WordDetailReq) (res *v1.WordDetailRes, err error) {
	w, err := daoWords.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询单词失败")
	}
	if w.IsEmpty() {
		return nil, errNotFound("未找到该单词")
	}
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	known := 0
	if cid > 0 {
		p, err := daoWordProg.Ctx(ctx).Where("word_id", req.ID).Where("child_id", cid).One()
		if err == nil && !p.IsEmpty() {
			known = p["known"].Int()
		}
	}
	return &v1.WordDetailRes{
		Word: v1.Word{
			ID:        w["id"].Int(),
			Level:     w["level"].Int(),
			Word:      w["word"].String(),
			Meaning:   w["meaning"].String(),
			Phonetic:  w["phonetic"].String(),
			Example:   w["example"].String(),
			CreatedAt: w["created_at"].String(),
		},
		Known: known,
	}, nil
}

// WordProgress 标记单词掌握状态；session_id 传入则走连击+场次计分。
func (s *sStudyPlanet) WordProgress(ctx context.Context, req *v1.WordProgressReq) (res *v1.WordProgressRes, err error) {
	known := 0
	if req.Known {
		known = 1
	}
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	if _, err := daoWordProg.Ctx(ctx).Data(doWordProg{
		WordId:       req.ID,
		ChildId:      cid,
		Known:        known,
		LastReviewed: gtime.Now(),
	}).Save(); err != nil {
		return nil, gerror.Wrap(err, "保存掌握状态失败")
	}
	res = &v1.WordProgressRes{Known: known}
	if req.SessionID > 0 {
		out, err := s.recordAnswer(ctx, cid, req.SessionID, req.ID, req.Known, 5, "",
			s.reviewRefs(ctx, cid, "words", []int{req.ID}))
		if err != nil {
			return nil, err
		}
		fillOutcome(&res.Combo, &res.BasePoints, &res.ComboBonus, &res.Review, &res.XP, &res.Correct, out)
		return res, nil
	}
	if req.Known {
		s.award(cid, 5, "单词认读:+5")
	}
	res.OK = true
	res.Correct = req.Known
	return res, nil
}

// ---------- 语文阅读 ----------

// ReadingDetail 阅读详情 + 题目列表。
func (s *sStudyPlanet) ReadingDetail(ctx context.Context, req *v1.ReadingDetailReq) (res *v1.ReadingDetailRes, err error) {
	rd, err := daoReadings.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询阅读失败")
	}
	if rd.IsEmpty() {
		return nil, errNotFound("未找到该阅读")
	}
	qs, err := daoReadingQ.Ctx(ctx).Where("reading_id", req.ID).Order("id").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询阅读题目失败")
	}
	res = &v1.ReadingDetailRes{
		Reading: v1.Reading{
			ID:      rd["id"].Int(),
			Title:   rd["title"].String(),
			Content: rd["content"].String(),
			Level:   rd["level"].Int(),
		},
		Questions: make([]v1.ReadingQuestion, 0, len(qs)),
	}
	for _, q := range qs {
		res.Questions = append(res.Questions, v1.ReadingQuestion{
			ID:        q["id"].Int(),
			ReadingID: q["reading_id"].Int(),
			Question:  q["question"].String(),
			OptionA:   q["option_a"].String(),
			OptionB:   q["option_b"].String(),
			OptionC:   q["option_c"].String(),
			OptionD:   q["option_d"].String(),
			Answer:    q["answer"].String(),
		})
	}
	return res, nil
}

// ReadingAnswer 阅读题目作答。
func (s *sStudyPlanet) ReadingAnswer(ctx context.Context, req *v1.ReadingAnswerReq) (res *v1.ReadingAnswerRes, err error) {
	q, err := daoReadingQ.Ctx(ctx).Where("id", req.QuestionID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询题目失败")
	}
	if q.IsEmpty() {
		return nil, errNotFound("未找到该题目")
	}
	answer := q["answer"].String()
	correct := strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(req.Answer))
	res = &v1.ReadingAnswerRes{Correct: correct, CorrectAnswer: answer}
	if req.SessionID > 0 {
		cid, err := s.resolveChild(ctx, req.StudentID)
		if err != nil {
			return nil, err
		}
		out, err := s.recordAnswer(ctx, cid, req.SessionID, req.QuestionID, correct, 2, answer,
			s.reviewRefs(ctx, cid, "reading", []int{req.QuestionID}))
		if err != nil {
			return nil, err
		}
		fillOutcome(&res.Combo, &res.BasePoints, &res.ComboBonus, &res.Review, &res.XP, nil, out)
		return res, nil
	}
	if correct {
		if cid, err := s.resolveChild(ctx, req.StudentID); err == nil && cid > 0 {
			s.award(cid, 2, "阅读答题:+2")
		}
	}
	return res, nil
}

// ---------- 数学题目 ----------

// ListMath 数学题目列表。
func (s *sStudyPlanet) ListMath(ctx context.Context, req *v1.ListMathReq) (res *v1.ListMathRes, err error) {
	m := daoMath.Ctx(ctx).Order("level", "id")
	if req.Level != "" {
		m = m.Where("level", req.Level)
	}
	all, err := m.All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询数学题失败")
	}
	out := make(v1.ListMathRes, 0, len(all))
	for _, r := range all {
		out = append(out, v1.MathProblem{
			ID:          r["id"].Int(),
			Level:       r["level"].Int(),
			Type:        r["type"].String(),
			Question:    r["question"].String(),
			Options:     r["options"].String(),
			Answer:      r["answer"].String(),
			Explanation: r["explanation"].String(),
		})
	}
	return &out, nil
}

// MathAnswer 数学题作答。
func (s *sStudyPlanet) MathAnswer(ctx context.Context, req *v1.MathAnswerReq) (res *v1.MathAnswerRes, err error) {
	p, err := daoMath.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询数学题失败")
	}
	if p.IsEmpty() {
		return nil, errNotFound("未找到该题目")
	}
	answer := p["answer"].String()
	correct := strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(req.Answer))
	res = &v1.MathAnswerRes{
		Correct:     correct,
		Explanation: p["explanation"].String(),
		Answer:      answer,
	}
	if req.SessionID > 0 {
		cid, err := s.resolveChild(ctx, req.StudentID)
		if err != nil {
			return nil, err
		}
		out, err := s.recordAnswer(ctx, cid, req.SessionID, req.ID, correct, 3, answer,
			s.reviewRefs(ctx, cid, "math", []int{req.ID}))
		if err != nil {
			return nil, err
		}
		fillOutcome(&res.Combo, &res.BasePoints, &res.ComboBonus, &res.Review, &res.XP, nil, out)
		return res, nil
	}
	if correct {
		if cid, err := s.resolveChild(ctx, req.StudentID); err == nil && cid > 0 {
			s.award(cid, 3, "数学答题:+3")
		}
	}
	return res, nil
}

// ---------- 每日任务 ----------

// ListTasks 学生任务列表（按 student_id，公开接口）。
func (s *sStudyPlanet) ListTasks(ctx context.Context, req *v1.ListTasksReq) (res *v1.ListTasksRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	m := daoTasks.Ctx(ctx).Where("child_id", cid)
	if req.Status != "" {
		m = m.Where("status", req.Status)
	}
	all, err := m.Order("due_date").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务失败")
	}
	today := todayStr()
	out := make(v1.ListTasksRes, 0, len(all))
	for _, r := range all {
		t := taskOf(r)
		if t.Status != "done" && t.DueDate != "" && t.DueDate < today {
			t.Status = "overdue"
		}
		out = append(out, t)
	}
	return &out, nil
}

// CompleteTask 完成任务（学生操作，公开接口）。
func (s *sStudyPlanet) CompleteTask(ctx context.Context, req *v1.CompleteTaskReq) (res *v1.CompleteTaskRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	t, err := daoTasks.Ctx(ctx).Where("id", req.ID).Where("child_id", cid).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务失败")
	}
	if t.IsEmpty() {
		return nil, errNotFound("未找到该任务")
	}
	if t["status"].String() == "done" {
		return &v1.CompleteTaskRes{OK: true, Already: true}, nil
	}
	if _, err := daoTasks.Ctx(ctx).Where("id", req.ID).Data(doTasks{
		Status:      "done",
		CompletedAt: gtime.Now(),
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "更新任务失败")
	}
	points := t["points"].Int()
	s.award(cid, points, "完成任务:+"+g.NewVar(points).String())
	return &v1.CompleteTaskRes{OK: true}, nil
}

// ---------- 积分 ----------

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

// ---------- 奖励 / 兑换 ----------

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

// ---------- 家长端 ----------

// ParentLogin 家长 PIN 登录（已废弃：多家长架构下 PIN 无法区分家长身份，强制走 Casdoor SSO）。
func (s *sStudyPlanet) ParentLogin(ctx context.Context, req *v1.ParentLoginReq) (res *v1.ParentLoginRes, err error) {
	return nil, errAuth("PIN 登录已停用，请使用 Casdoor SSO 登录")
}

// ---------- 以下为家长鉴权后接口 ----------

// AddTask 发布任务（家长鉴权）：学生必须归属当前家长。
func (s *sStudyPlanet) AddTask(ctx context.Context, req *v1.AddTaskReq) (res *v1.AddTaskRes, err error) {
	if strings.TrimSpace(req.Title) == "" {
		return nil, errParam("请填写任务名称")
	}
	if err := s.ensureChildOwned(ctx, req.StudentID); err != nil {
		return nil, err
	}
	data := doTasks{
		Title:   req.Title,
		Type:    req.Type,
		Points:  req.Points,
		Status:  "pending",
		ChildId: req.StudentID,
	}
	// due_date 为空时留空（MySQL DATE 列不接受空字符串）
	if strings.TrimSpace(req.DueDate) != "" {
		if t, err := gtime.StrToTime(req.DueDate, "Y-m-d"); err == nil {
			data.DueDate = t
		} else {
			data.DueDate = gtime.NewFromStr(req.DueDate)
		}
	}
	if _, err := daoTasks.Ctx(ctx).Data(data).Insert(); err != nil {
		return nil, gerror.Wrap(err, "发布任务失败")
	}
	return &v1.AddTaskRes{OK: true}, nil
}

// DeleteTask 删除任务（家长鉴权）：任务须属于当前家长的学生。
func (s *sStudyPlanet) DeleteTask(ctx context.Context, req *v1.DeleteTaskReq) (res *v1.DeleteTaskRes, err error) {
	t, err := daoTasks.Ctx(ctx).Fields("id,child_id").Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询任务失败")
	}
	if t.IsEmpty() {
		return nil, errNotFound("未找到该任务")
	}
	if err := s.ensureChildOwned(ctx, t["child_id"].Int()); err != nil {
		return nil, err
	}
	if _, err := daoTasks.Ctx(ctx).Where("id", req.ID).Delete(); err != nil {
		return nil, gerror.Wrap(err, "删除任务失败")
	}
	return &v1.DeleteTaskRes{OK: true}, nil
}

// AddReward 添加奖励（家长鉴权）：奖励归属当前家长，仅自家孩子可见可兑。
func (s *sStudyPlanet) AddReward(ctx context.Context, req *v1.AddRewardReq) (res *v1.AddRewardRes, err error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errParam("请填写奖励名称")
	}
	parentID := ctxParentID(ctx)
	if parentID <= 0 {
		return nil, errAuth("登录状态缺少家长身份，请重新登录")
	}
	if _, err := daoRewards.Ctx(ctx).Data(doRewards{
		Name:       req.Name,
		CostPoints: req.CostPoints,
		Status:     "active",
		ParentId:   parentID,
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

// SetPin 修改家长 PIN（已废弃：PIN 登录停用，保留接口返回提示避免旧前端报 404）。
func (s *sStudyPlanet) SetPin(ctx context.Context, req *v1.SetPinReq) (res *v1.SetPinRes, err error) {
	return nil, errParam("PIN 登录已停用，无需设置 PIN")
}

// ---------- 学生账号管理（家长鉴权后） ----------

// ListStudents 学生列表（挂可选鉴权）：
//   - 匿名（学生端未登录家长）：返回空列表——匿名用户不应看到任何家庭的孩子；
//     前端家长登录（Casdoor）后本接口自动返回本家孩子。
//   - 带家长 token：只返回归属自己的孩子（数据隔离）。
func (s *sStudyPlanet) ListStudents(ctx context.Context, req *v1.ListStudentsReq) (res *v1.ListStudentsRes, err error) {
	parentID := ctxParentID(ctx)
	out := make(v1.ListStudentsRes, 0, 4)
	if parentID <= 0 {
		return &out, nil
	}
	all, err := daoChildren.Ctx(ctx).Where("parent_id", parentID).Order("id").All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询学生失败")
	}
	for _, r := range all {
		out = append(out, studentOf(r))
	}
	return &out, nil
}

// CreateStudent 新建学生账号：自动归属当前登录家长。
func (s *sStudyPlanet) CreateStudent(ctx context.Context, req *v1.CreateStudentReq) (res *v1.CreateStudentRes, err error) {
	parentID := ctxParentID(ctx)
	if parentID <= 0 {
		return nil, errAuth("登录状态缺少家长身份，请重新登录")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errParam("请填写学生姓名")
	}
	avatar, grade := req.Avatar, req.Grade
	if avatar == "" {
		avatar = "🚀"
	}
	if grade <= 0 {
		grade = 5
	}
	if req.Username != "" {
		cnt, err := daoChildren.Ctx(ctx).Where("username", req.Username).Count()
		if err != nil {
			return nil, gerror.Wrap(err, "校验用户名失败")
		}
		if cnt > 0 {
			return nil, errParam("用户名已被使用")
		}
	}
	id, err := daoChildren.Ctx(ctx).Data(doChildren{
		Name:     req.Name,
		Username: req.Username,
		Avatar:   avatar,
		Grade:    grade,
		ParentId: parentID,
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.Wrap(err, "创建学生失败")
	}
	rec, err := daoChildren.Ctx(ctx).Where("id", id).One()
	if err != nil || rec.IsEmpty() {
		return nil, errNotFound("查询新学生失败")
	}
	st := v1.CreateStudentRes(studentOf(rec))
	return &st, nil
}

// UpdateStudent 修改学生信息（姓名/头像/年级/用户名）：仅限自己的孩子。
func (s *sStudyPlanet) UpdateStudent(ctx context.Context, req *v1.UpdateStudentReq) (res *v1.UpdateStudentRes, err error) {
	if err := s.ensureChildOwned(ctx, req.ID); err != nil {
		return nil, err
	}
	if req.Name != nil && *req.Name == "" {
		return nil, errParam("姓名不能为空")
	}
	if req.Username != nil && *req.Username != "" {
		cnt, err := daoChildren.Ctx(ctx).Where("username", *req.Username).Where("id<>?", req.ID).Count()
		if err != nil {
			return nil, gerror.Wrap(err, "校验用户名失败")
		}
		if cnt > 0 {
			return nil, errParam("用户名已被使用")
		}
	}
	data := g.Map{}
	if req.Name != nil {
		data["name"] = *req.Name
	}
	if req.Username != nil {
		data["username"] = *req.Username
	}
	if req.Avatar != nil {
		data["avatar"] = *req.Avatar
	}
	if req.Grade != nil {
		data["grade"] = *req.Grade
	}
	if len(data) > 0 {
		if _, err := daoChildren.Ctx(ctx).Where("id", req.ID).Data(data).Update(); err != nil {
			return nil, gerror.Wrap(err, "更新学生失败")
		}
	}
	rec, err := daoChildren.Ctx(ctx).Where("id", req.ID).One()
	if err != nil || rec.IsEmpty() {
		return nil, errNotFound("学生不存在")
	}
	st := v1.UpdateStudentRes(studentOf(rec))
	return &st, nil
}

// DeleteStudent 删除学生；至少保留一个，且清空其学习数据。仅限自己的孩子。
func (s *sStudyPlanet) DeleteStudent(ctx context.Context, req *v1.DeleteStudentReq) (res *v1.DeleteStudentRes, err error) {
	if err := s.ensureChildOwned(ctx, req.ID); err != nil {
		return nil, err
	}
	cnt, err := daoChildren.Ctx(ctx).Where("parent_id", ctxParentID(ctx)).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "统计学生失败")
	}
	if cnt <= 1 {
		return nil, errParam("至少保留一个学生账号")
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先清场次答题（依赖 session_id），再清其余学习数据，最后删学生
		if _, err := tx.Exec("DELETE FROM session_answers WHERE session_id IN (SELECT id FROM practice_sessions WHERE child_id=?)", req.ID); err != nil {
			return err
		}
		for _, d := range []*gdb.Model{
			daoWrongQ.Ctx(ctx).Where("child_id", req.ID),
			daoSessions.Ctx(ctx).Where("child_id", req.ID),
			daoWordProg.Ctx(ctx).Where("child_id", req.ID),
			daoPointsLog.Ctx(ctx).Where("child_id", req.ID),
			daoTasks.Ctx(ctx).Where("child_id", req.ID),
			daoRedempt.Ctx(ctx).Where("child_id", req.ID),
			daoChildren.Ctx(ctx).Where("id", req.ID),
		} {
			if _, err := d.Delete(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, gerror.Wrap(err, "删除学生失败")
	}
	return &v1.DeleteStudentRes{OK: true}, nil
}
