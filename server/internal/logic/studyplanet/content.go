package studyplanet

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "studyplanet/api/studyplanet/v1"
	"studyplanet/internal/contentlib"
)

// ---------- 动态内容库：出题 / 判分 / 导入（学习内容全部来自数据库） ----------

// ListSubjects 学科目录：GET /api/subjects?grade=5
// 返回学科列表 + 每科题量，前端动态渲染学习地图。
// 传 grade 时只返回该学段开设的学科（如 5 年级不出物理/化学）。
func (s *sStudyPlanet) ListSubjects(ctx context.Context, req *v1.ListSubjectsReq) (res *v1.ListSubjectsRes, err error) {
	ss, err := contentlib.ListSubjects(ctx, g.DB())
	if err != nil {
		return nil, gerror.Wrap(err, "查询学科失败")
	}
	if req.Grade >= 1 && req.Grade <= 9 {
		filtered := make([]contentlib.Subject, 0, len(ss))
		for _, sub := range ss {
			if req.Grade >= sub.MinGrade && req.Grade <= sub.MaxGrade {
				filtered = append(filtered, sub)
			}
		}
		ss = filtered
	}
	res = &v1.ListSubjectsRes{}
	for _, sub := range ss {
		*res = append(*res, subjectOf(sub, 0))
	}
	counts, err := contentlib.CountBySubject(ctx, g.DB())
	if err == nil {
		for i := range *res {
			(*res)[i].Count = counts[(*res)[i].Code]
		}
	}
	return res, nil
}

// subjectOf contentlib 学科 → api 学科条目。
func subjectOf(sub contentlib.Subject, count int) v1.Subject {
	return v1.Subject{
		Code:     sub.Code,
		Name:     sub.Name,
		Icon:     sub.Icon,
		Color:    sub.Color,
		MinGrade: sub.MinGrade,
		MaxGrade: sub.MaxGrade,
		Sort:     sub.Sort,
		Enabled:  sub.Enabled,
		Count:    count,
	}
}

// pubQuestionOf questions 记录 → 对外发布的题目（不含 answer，防泄露）。
func pubQuestionOf(id, grade, difficulty int, subject, topic, qtype, passage, question, optionsJSON string) v1.PubQuestion {
	var opts []string
	_ = json.Unmarshal([]byte(optionsJSON), &opts)
	if opts == nil {
		opts = []string{}
	}
	return v1.PubQuestion{
		ID:         id,
		Subject:    subject,
		Grade:      grade,
		Topic:      topic,
		QType:      qtype,
		Passage:    passage,
		Question:   question,
		Options:    opts,
		Difficulty: difficulty,
	}
}

// PickQuestions 从内容库随机抽题：GET /api/content/pick?subject=math&grade=5&limit=5
// 返回的题目不带 answer（不泄露答案给前端），前端作答后走 /api/content/answer 判分。
func (s *sStudyPlanet) PickQuestions(ctx context.Context, req *v1.PickQuestionsReq) (res *v1.PickQuestionsRes, err error) {
	if strings.TrimSpace(req.Subject) == "" {
		return nil, errParam("subject 不能为空")
	}
	limit := req.Limit
	if limit <= 0 || limit > 20 {
		limit = 5
	}
	m := daoQuestions.Ctx(ctx).Where("subject", req.Subject).Where("enabled", 1)
	if req.Grade > 0 {
		// 学生年级：取该年级及相邻难度（grade-1..grade+1），保证题量充足
		m = m.WhereGTE("grade", req.Grade-1).WhereLTE("grade", req.Grade+1)
	}
	all, err := m.OrderRandom().Limit(limit).All()
	if err != nil {
		return nil, gerror.Wrap(err, "抽题失败")
	}
	out := make(v1.PickQuestionsRes, 0, len(all))
	for _, r := range all {
		out = append(out, pubQuestionOf(r["id"].Int(), r["grade"].Int(), r["difficulty"].Int(),
			r["subject"].String(), r["topic"].String(), r["qtype"].String(),
			r["passage"].String(), r["question"].String(), r["options"].String(),
		))
	}
	return &out, nil
}

// ContentAnswer 统一判分：POST /api/content/answer
// body: {id, answer, session_id?}；session_id 传入则复用连击+XP+错题本链路。
func (s *sStudyPlanet) ContentAnswer(ctx context.Context, req *v1.ContentAnswerReq) (res *v1.ContentAnswerRes, err error) {
	row, err := daoQuestions.Ctx(ctx).Where("id", req.ID).Where("enabled", 1).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询题目失败")
	}
	if row.IsEmpty() {
		return nil, errNotFound("题目不存在")
	}
	answer := row["answer"].String()
	correct := strings.EqualFold(strings.TrimSpace(answer), strings.TrimSpace(req.Answer))
	res = &v1.ContentAnswerRes{
		Correct:     correct,
		Answer:      answer,
		Explanation: row["explanation"].String(),
	}
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if req.SessionID > 0 && cid > 0 {
		// 复用多邻国链路：连击 / XP / 错题本（错题统一记录）
		out, err := s.recordAnswer(ctx, cid, req.SessionID, req.ID, correct, 3, answer,
			s.reviewRefs(ctx, cid, row["subject"].String(), []int{req.ID}))
		if err != nil {
			return nil, err
		}
		fillOutcome(&res.Combo, &res.BasePoints, &res.ComboBonus, &res.Review, &res.XP, nil, out)
		return res, nil
	}
	if correct && cid > 0 {
		s.award(cid, 3, "内容库答题:+3")
	}
	return res, nil
}

// ContentItem 按 id 取单题（不含答案）：GET /api/content/item?id=
// 错题本巩固复习时回取题目内容用。
func (s *sStudyPlanet) ContentItem(ctx context.Context, req *v1.ContentItemReq) (res *v1.ContentItemRes, err error) {
	row, err := daoQuestions.Ctx(ctx).Where("id", req.ID).Where("enabled", 1).One()
	if err != nil {
		return nil, gerror.Wrap(err, "查询题目失败")
	}
	if row.IsEmpty() {
		return nil, errNotFound("题目不存在")
	}
	item := v1.ContentItemRes(pubQuestionOf(row["id"].Int(), row["grade"].Int(), row["difficulty"].Int(),
		row["subject"].String(), row["topic"].String(), row["qtype"].String(),
		row["passage"].String(), row["question"].String(), row["options"].String(),
	))
	return &item, nil
}

// ImportContent 通用题目导入（家长身份）：POST /api/parent/content/import
// body: {"questions": [{subject,grade,topic,qtype,passage?,question,options[],answer,explanation?,difficulty?,source?}]}
// 按 content_hash 去重，重复导入自动跳过——以后采集新资料直接调本接口，无需改源码。
func (s *sStudyPlanet) ImportContent(ctx context.Context, req *v1.ImportContentReq) (res *v1.ImportContentRes, err error) {
	if len(req.Questions) == 0 {
		return nil, errParam("questions 不能为空")
	}
	if len(req.Questions) > 2000 {
		return nil, errParam("单次最多导入 2000 题，请分批")
	}
	qs := make([]contentlib.Question, 0, len(req.Questions))
	for _, q := range req.Questions {
		qs = append(qs, contentlib.Question{
			Subject:     q.Subject,
			Grade:       q.Grade,
			Topic:       q.Topic,
			QType:       q.QType,
			Passage:     q.Passage,
			Question:    q.Question,
			Options:     q.Options,
			Answer:      q.Answer,
			Explanation: q.Explanation,
			Difficulty:  q.Difficulty,
			Source:      q.Source,
		})
	}
	imported, skipped, err := contentlib.ImportQuestions(ctx, g.DB(), qs)
	if err != nil {
		return nil, gerror.Wrap(err, "导入题目失败")
	}
	return &v1.ImportContentRes{Imported: imported, Skipped: skipped, Total: len(req.Questions)}, nil
}

// SubjectStats 内容库统计（家长端）：GET /api/parent/content/stats
func (s *sStudyPlanet) SubjectStats(ctx context.Context, req *v1.SubjectStatsReq) (res *v1.SubjectStatsRes, err error) {
	ss, err := contentlib.ListSubjects(ctx, g.DB())
	if err != nil {
		return nil, gerror.Wrap(err, "查询学科失败")
	}
	counts, err := contentlib.CountBySubject(ctx, g.DB())
	if err != nil {
		return nil, gerror.Wrap(err, "统计题量失败")
	}
	total := 0
	for _, c := range counts {
		total += c
	}
	res = &v1.SubjectStatsRes{Total: total, Subjects: make([]v1.Subject, 0, len(ss))}
	for _, sub := range ss {
		res.Subjects = append(res.Subjects, subjectOf(sub, counts[sub.Code]))
	}
	return res, nil
}
