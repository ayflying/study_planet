// Package studyplanet 语文阅读模块：阅读详情 / 题目作答。
package studyplanet

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "studyplanet/api/studyplanet/v1"
)

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
