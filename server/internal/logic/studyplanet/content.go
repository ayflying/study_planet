package studyplanet

import (
	"encoding/json"
	"math/rand"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"studyplanet/internal/contentlib"
)

// ---------- 动态内容库：出题 / 判分 / 导入（学习内容全部来自数据库） ----------

// pickToken 出题令牌：服务端把题目缓存在 questions 表外（按 token 不可行，改为
// 无状态方案——前端每次取题集合，判分时只校验答案是否正确，题目 id 回传即可）。

// ListSubjects 学科目录：GET /api/subjects?grade=5
// 返回学科列表 + 每科题量，前端动态渲染学习地图。
// 传 grade 时只返回该学段开设的学科（如 5 年级不出物理/化学）。
func (s *sStudyPlanet) ListSubjects(r *ghttp.Request) {
	ss, err := contentlib.ListSubjects(s.DB)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	if g := r.GetQuery("grade").Int(); g >= 1 && g <= 9 {
		filtered := make([]contentlib.Subject, 0, len(ss))
		for _, sub := range ss {
			if g >= sub.MinGrade && g <= sub.MaxGrade {
				filtered = append(filtered, sub)
			}
		}
		ss = filtered
	}
	counts, _ := contentlib.CountBySubject(s.DB)
	type item struct {
		contentlib.Subject
		Count int `json:"count"`
	}
	items := make([]item, 0, len(ss))
	for _, sub := range ss {
		items = append(items, item{Subject: sub, Count: counts[sub.Code]})
	}
	s.ok(r, items)
}

// PickQuestions 从内容库随机抽题：GET /api/content/pick?subject=math&grade=5&limit=5
// 返回的题目不带 answer（不泄露答案给前端），前端作答后走 /api/content/answer 判分。
func (s *sStudyPlanet) PickQuestions(r *ghttp.Request) {
	subject := strings.TrimSpace(r.GetQuery("subject").String())
	if subject == "" {
		s.fail(r, 400, "subject 不能为空")
		return
	}
	grade := r.GetQuery("grade").Int()
	limit := r.GetQuery("limit").Int()
	if limit <= 0 || limit > 20 {
		limit = 5
	}
	q := "SELECT id,subject,grade,topic,qtype,passage,question,options,explanation,difficulty FROM questions WHERE subject=? AND enabled=1"
	args := []interface{}{subject}
	if grade > 0 {
		// 学生年级：取该年级及相邻难度（grade-1..grade+1），保证题量充足
		q += " AND grade BETWEEN ? AND ?"
		args = append(args, grade-1, grade+1)
	}
	q += " ORDER BY RANDOM() LIMIT ?"
	if s.DB.DriverName() == "mysql" {
		q = strings.Replace(q, "ORDER BY RANDOM()", "ORDER BY RAND()", 1)
	}
	args = append(args, limit)
	var rows []struct {
		ID          int    `db:"id"`
		Subject     string `db:"subject"`
		Grade       int    `db:"grade"`
		Topic       string `db:"topic"`
		QType       string `db:"qtype"`
		Passage     *string `db:"passage"`
		Question    string `db:"question"`
		OptionsJSON string `db:"options"`
		Explanation *string `db:"explanation"`
		Difficulty  int    `db:"difficulty"`
	}
	if err := s.DB.Select(&rows, q, args...); err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	type pubQuestion struct {
		ID          int      `json:"id"`
		Subject     string   `json:"subject"`
		Grade       int      `json:"grade"`
		Topic       string   `json:"topic"`
		QType       string   `json:"qtype"`
		Passage     string   `json:"passage,omitempty"`
		Question    string   `json:"question"`
		Options     []string `json:"options"`
		Difficulty  int      `json:"difficulty"`
	}
	out := make([]pubQuestion, 0, len(rows))
	for _, row := range rows {
		var opts []string
		_ = json.Unmarshal([]byte(row.OptionsJSON), &opts)
		var passage string
		if row.Passage != nil {
			passage = *row.Passage
		}
		pq := pubQuestion{
			ID: row.ID, Subject: row.Subject, Grade: row.Grade, Topic: row.Topic,
			QType: row.QType, Passage: passage, Question: row.Question,
			Options: opts, Difficulty: row.Difficulty,
		}
		if pq.Options == nil {
			pq.Options = []string{}
		}
		out = append(out, pq)
	}
	s.ok(r, out)
}

// ContentAnswer 统一判分：POST /api/content/answer
// body: {id, answer, session_id?}；session_id 传入则复用连击+XP+错题本链路。
func (s *sStudyPlanet) ContentAnswer(r *ghttp.Request) {
	var body struct {
		ID        int    `json:"id"`
		Answer    string `json:"answer"`
		SessionID int    `json:"session_id"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误")
		return
	}
	var row struct {
		ID       int    `db:"id"`
		Subject  string `db:"subject"`
		Answer   string `db:"answer"`
		Explain  *string `db:"explanation"`
	}
	if err := s.DB.Get(&row, "SELECT id,subject,answer,explanation FROM questions WHERE id=? AND enabled=1", body.ID); err != nil {
		s.fail(r, 404, "题目不存在")
		return
	}
	correct := strings.EqualFold(strings.TrimSpace(row.Answer), strings.TrimSpace(body.Answer))
	if body.SessionID > 0 {
		// 复用现有多邻国链路：连击 / XP / 错题本（subject 与旧链路同名映射，错题统一记录）
		s.recordAnswer(r, body.SessionID, body.ID, correct, 3, row.Answer, s.reviewRefs(r, row.Subject, []int{body.ID}))
		return
	}
	if correct {
		if cid := s.resolveChild(r); cid > 0 {
			s.award(cid, 3, "内容库答题:+3")
		}
	}
	explanation := ""
	if row.Explain != nil {
		explanation = *row.Explain
	}
	s.ok(r, map[string]interface{}{"correct": correct, "answer": row.Answer, "explanation": explanation})
}

var _ = rand.Intn // 占位避免未使用告警

// ContentItem 按 id 取单题（不含答案）：GET /api/content/item?id=
// 错题本巩固复习时回取题目内容用。
func (s *sStudyPlanet) ContentItem(r *ghttp.Request) {
	id := s.idParam(r)
	var row struct {
		ID          int     `db:"id"`
		Subject     string  `db:"subject"`
		Grade       int     `db:"grade"`
		Topic       string  `db:"topic"`
		QType       string  `db:"qtype"`
		Passage     *string `db:"passage"`
		Question    string  `db:"question"`
		OptionsJSON string  `db:"options"`
		Difficulty  int     `db:"difficulty"`
	}
	if err := s.DB.Get(&row, "SELECT id,subject,grade,topic,qtype,passage,question,options,difficulty FROM questions WHERE id=? AND enabled=1", id); err != nil {
		s.fail(r, 404, "题目不存在")
		return
	}
	var opts []string
	_ = json.Unmarshal([]byte(row.OptionsJSON), &opts)
	if opts == nil {
		opts = []string{}
	}
	passage := ""
	if row.Passage != nil {
		passage = *row.Passage
	}
	s.ok(r, map[string]interface{}{
		"id": row.ID, "subject": row.Subject, "grade": row.Grade, "topic": row.Topic,
		"qtype": row.QType, "passage": passage, "question": row.Question,
		"options": opts, "difficulty": row.Difficulty,
	})
}

// ImportContent 通用题目导入（家长身份）：POST /api/parent/content/import
// body: {"questions": [{subject,grade,topic,qtype,passage?,question,options[],answer,explanation?,difficulty?,source?}]}
// 按 content_hash 去重，重复导入自动跳过——以后采集新资料直接调本接口，无需改源码。
func (s *sStudyPlanet) ImportContent(r *ghttp.Request) {
	var body struct {
		Questions []contentlib.Question `json:"questions"`
	}
	if err := r.Parse(&body); err != nil {
		s.fail(r, 400, "请求格式错误: "+err.Error())
		return
	}
	if len(body.Questions) == 0 {
		s.fail(r, 400, "questions 不能为空")
		return
	}
	if len(body.Questions) > 2000 {
		s.fail(r, 400, "单次最多导入 2000 题，请分批")
		return
	}
	imported, skipped, err := contentlib.ImportQuestions(s.DB, body.Questions)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	s.ok(r, map[string]interface{}{"imported": imported, "skipped": skipped, "total": len(body.Questions)})
}

// SubjectStats 内容库统计（家长端）：GET /api/parent/content/stats
func (s *sStudyPlanet) SubjectStats(r *ghttp.Request) {
	ss, err := contentlib.ListSubjects(s.DB)
	if err != nil {
		s.fail(r, 500, err.Error())
		return
	}
	counts, _ := contentlib.CountBySubject(s.DB)
	var total int
	for _, c := range counts {
		total += c
	}
	type item struct {
		contentlib.Subject
		Count int `json:"count"`
	}
	items := make([]item, 0, len(ss))
	for _, sub := range ss {
		items = append(items, item{Subject: sub, Count: counts[sub.Code]})
	}
	s.ok(r, map[string]interface{}{"total": total, "subjects": items})
}
