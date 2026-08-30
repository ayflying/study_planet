// Package studyplanet 数学题目模块：列表 / 作答。
package studyplanet

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "studyplanet/api/studyplanet/v1"
)

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
