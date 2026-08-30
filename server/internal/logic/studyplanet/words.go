// Package studyplanet 单词卡片模块：列表 / 详情 / 掌握进度。
package studyplanet

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

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
