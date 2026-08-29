package studyplanet

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
	"studyplanet/internal/leaderboard"
)

// ---------- 错题本：答错登记 / 巩固练习出题 ----------

// recordWrong 答错时登记错题（已 resolved 的重新激活，wrong_count 累加）。
// 只做附加记录，失败不影响作答主流程。
// 说明：GF 的 OnDuplicate map 值按「列名引用」处理，累加表达式需用 gdb.Raw 显式更新。
func (s *sStudyPlanet) recordWrong(ctx context.Context, childID int, subject string, refID int) {
	where := func(m *gdb.Model) *gdb.Model {
		return m.Where("child_id", childID).Where("subject", subject).Where("ref_id", refID)
	}
	cnt, err := where(daoWrongQ.Ctx(ctx)).Count()
	if err == nil && cnt > 0 {
		if _, err := where(daoWrongQ.Ctx(ctx)).Data(g.Map{
			"wrong_count":   gdb.Raw("wrong_count+1"),
			"resolved":      0,
			"last_wrong_at": nowStr(),
		}).Update(); err != nil {
			gLog("错题累加失败 child=%d %s#%d: %v", childID, subject, refID, err)
		}
		return
	}
	if _, err := daoWrongQ.Ctx(ctx).Data(doWrongQ{
		ChildId:     childID,
		Subject:     subject,
		RefId:       refID,
		WrongCount:  1,
		Resolved:    0,
		LastWrongAt: gtime.Now(),
	}).Insert(); err != nil {
		gLog("错题登记失败 child=%d %s#%d: %v", childID, subject, refID, err)
	}
}

// resolveWrong 巩固练习中答对后消除错题。
func (s *sStudyPlanet) resolveWrong(ctx context.Context, childID int, subject string, refID int) {
	if _, err := daoWrongQ.Ctx(ctx).
		Where("child_id", childID).
		Where("subject", subject).
		Where("ref_id", refID).
		Data(doWrongQ{Resolved: 1, LastReviewedAt: gtime.Now()}).Update(); err != nil {
		gLog("错题消除失败 child=%d %s#%d: %v", childID, subject, refID, err)
	}
}

// wrongIDs 学生当前未消除的错题 id 列表（按错误次数倒序，优先巩固错得多的）。
func (s *sStudyPlanet) wrongIDs(ctx context.Context, childID int, subject string, limit int) []int {
	if limit <= 0 {
		return nil
	}
	all, err := daoWrongQ.Ctx(ctx).
		Fields("ref_id").
		Where("child_id", childID).
		Where("subject", subject).
		Where("resolved", 0).
		OrderDesc("wrong_count").
		OrderDesc("last_wrong_at").
		Limit(limit).All()
	if err != nil {
		return nil
	}
	ids := make([]int, 0, len(all))
	for _, r := range all {
		ids = append(ids, r["ref_id"].Int())
	}
	return ids
}

// isReviewRef 判断本题是否来自错题本（用于响应标记与答对消除）。
func isReviewRef(reviewRefs map[int]bool, id int) bool { return reviewRefs[id] }

// reviewRefs 构建本题集合中的错题标记（服务端自行判断，避免信任前端）。
func (s *sStudyPlanet) reviewRefs(ctx context.Context, childID int, subject string, refIDs []int) map[int]bool {
	if childID < 0 || len(refIDs) == 0 {
		return nil
	}
	all, err := daoWrongQ.Ctx(ctx).
		Fields("ref_id").
		Where("child_id", childID).
		Where("subject", subject).
		Where("resolved", 0).
		Where("ref_id", refIDs).All()
	if err != nil {
		return nil
	}
	m := make(map[int]bool, len(all))
	for _, r := range all {
		m[r["ref_id"].Int()] = true
	}
	return m
}

// ListWrongQuestions 学生错题本（可选按 subject 过滤）。
func (s *sStudyPlanet) ListWrongQuestions(ctx context.Context, req *v1.ListWrongQuestionsReq) (res *v1.ListWrongQuestionsRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid < 0 {
		return nil, errNotFound("学生不存在")
	}
	m := daoWrongQ.Ctx(ctx).Where("child_id", cid)
	if req.Subject != "" {
		m = m.Where("subject", req.Subject)
	}
	all, err := m.Order("resolved").OrderDesc("wrong_count").OrderDesc("last_wrong_at").Limit(200).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询错题本失败")
	}
	out := make(v1.ListWrongQuestionsRes, 0, len(all))
	for _, r := range all {
		out = append(out, v1.WrongQuestion{
			ID:             r["id"].Int(),
			ChildID:        r["child_id"].Int(),
			Subject:        r["subject"].String(),
			RefID:          r["ref_id"].Int(),
			WrongCount:     r["wrong_count"].Int(),
			Resolved:       r["resolved"].Int(),
			LastWrongAt:    r["last_wrong_at"].String(),
			LastReviewedAt: r["last_reviewed_at"].String(),
		})
	}
	return &out, nil
}

// ---------- 每周经验排行榜 ----------

// WeeklyLeaderboard 周榜：返回当前 ISO 周经验值最高的学生名单。
func (s *sStudyPlanet) WeeklyLeaderboard(ctx context.Context, req *v1.WeeklyLeaderboardReq) (res *v1.WeeklyLeaderboardRes, err error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	week := leaderboard.WeekKey(time.Now())
	var entries []leaderboard.Entry
	if s.Board != nil {
		entries = s.Board.Top(ctx, week, limit, func(id int) (string, string) {
			rec, err := daoChildren.Ctx(ctx).Fields("name,avatar").Where("id", id).One()
			if err != nil || rec.IsEmpty() {
				return "", ""
			}
			return rec["name"].String(), rec["avatar"].String()
		})
	}
	// 学生自己的排名（可能不在前 N）
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	myXP, myRank := 0, 0
	if cid > 0 {
		v, err := daoChildren.Ctx(ctx).Fields("COALESCE(xp,0) AS xp").Where("id", cid).Value()
		if err == nil {
			myXP = v.Int()
		}
		for _, e := range entries {
			if e.ChildID == cid {
				myRank = e.Rank
				break
			}
		}
	}
	res = &v1.WeeklyLeaderboardRes{
		Week:    week,
		Redis:   s.Board != nil && s.Board.Enabled(),
		Entries: make([]v1.LeaderboardEntry, 0, len(entries)),
		MyXP:    myXP,
		MyRank:  myRank,
		MyID:    cid,
	}
	for _, e := range entries {
		res.Entries = append(res.Entries, v1.LeaderboardEntry{
			Rank:    e.Rank,
			ChildID: e.ChildID,
			Name:    e.Name,
			Avatar:  e.Avatar,
			XP:      e.XP,
		})
	}
	return res, nil
}
