// Package studyplanet 对战查询接口：段位榜 / 历史战绩，以及孩子信息容错读取。
package studyplanet

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "studyplanet/api/studyplanet/v1"
	"studyplanet/internal/battle"
)

// tierOf 奖杯数 → 段位：直接复用 battle 包的唯一实现，避免两处规则漂移。
func tierOf(trophies int) (string, string) {
	return battle.TierName(trophies)
}

// BattleRank 对战段位榜：按奖杯数排序，附我的段位。
func (s *sStudyPlanet) BattleRank(ctx context.Context, req *v1.BattleRankReq) (res *v1.BattleRankRes, err error) {
	limit := req.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	rows, err := g.DB().Model("battle_scores b").Ctx(ctx).
		LeftJoin("children c", "c.id=b.child_id").
		Fields("b.child_id", "b.trophies", "b.wins", "b.battles", "c.name", "c.avatar").
		Order("b.trophies DESC").Order("b.wins DESC").Limit(limit).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询段位榜失败")
	}
	res = &v1.BattleRankRes{}
	res.Entries = make([]v1.BattleRankEntry, 0, len(rows))
	for i, r := range rows {
		tier, emoji := tierOf(r["trophies"].Int())
		res.Entries = append(res.Entries, v1.BattleRankEntry{
			Rank: i + 1, ChildID: r["child_id"].Int(),
			Name: r["name"].String(), Avatar: r["avatar"].String(),
			Trophies: r["trophies"].Int(), Tier: tier, TierEmoji: emoji,
			Wins: r["wins"].Int(), Battles: r["battles"].Int(),
		})
	}
	// 我的段位
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid > 0 {
		my, err := g.DB().Model("battle_scores").Ctx(ctx).Where("child_id", cid).One()
		if err != nil {
			return nil, gerror.Wrap(err, "查询我的段位失败")
		}
		if !my.IsEmpty() {
			tier, emoji := tierOf(my["trophies"].Int())
			myTrophies := my["trophies"].Int()
			rank := 1
			for _, e := range res.Entries {
				if e.Trophies > myTrophies {
					rank++
				}
			}
			cnt, _ := g.DB().Model("battle_scores").Ctx(ctx).Count()
			res.My = &v1.BattleRankEntry{
				Rank: rank, ChildID: cid,
				Name: s.childName(ctx, cid), Avatar: s.childAvatar(ctx, cid),
				Trophies: myTrophies, Tier: tier, TierEmoji: emoji,
				Wins: my["wins"].Int(), Battles: my["battles"].Int(),
			}
			if res.My.Rank == 1 && cnt > 1 {
				// 榜内没有比自己分高的（自己不在榜内时）
				if len(res.Entries) > 0 && res.Entries[0].ChildID != cid && res.Entries[0].Trophies >= myTrophies {
					res.My.Rank = res.Entries[0].Rank + 1
				}
			}
		}
	}
	return res, nil
}

// BattleHistory 我的历史对战（最近 30 条）。
func (s *sStudyPlanet) BattleHistory(ctx context.Context, req *v1.BattleHistoryReq) (res *v1.BattleHistoryRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid <= 0 {
		empty := v1.BattleHistoryRes{}
		return &empty, nil
	}
	rows, err := g.DB().Model("battles").Ctx(ctx).
		Where("p1_id", cid).WhereOr("p2_id", cid).
		Where("status", "finished").
		Order("id DESC").Limit(30).All()
	if err != nil {
		return nil, gerror.Wrap(err, "查询对战历史失败")
	}
	out := make(v1.BattleHistoryRes, 0, len(rows))
	for _, r := range rows {
		p1 := r["p1_id"].Int()
		myScore, oppScore := r["p1_score"].Int(), r["p2_score"].Int()
		oppID := r["p2_id"].Int()
		isP1 := p1 == cid
		if !isP1 {
			myScore, oppScore = oppScore, myScore
			oppID = p1
		}
		result := "draw"
		if myScore > oppScore {
			result = "win"
		} else if myScore < oppScore {
			result = "lose"
		}
		trophies := r["trophies1"].Int()
		if !isP1 {
			trophies = r["trophies2"].Int()
		}
		opName, opAvatar := "练习机器人", "🤖"
		if r["p2_robot"].Int() == 0 && oppID > 0 {
			opName = s.childName(ctx, oppID)
			opAvatar = s.childAvatar(ctx, oppID)
		}
		created := r["created_at"].Time()
		createdStr := ""
		if !created.IsZero() && created.Year() > 2000 {
			createdStr = created.Format("01-02 15:04")
		}
		out = append(out, v1.BattleHistoryEntry{
			ID: r["id"].Int(), Opponent: opName, OpponentAvatar: opAvatar,
			MyScore: myScore, OppScore: oppScore, Result: result,
			Trophies: trophies, CreatedAt: createdStr,
		})
	}
	return &out, nil
}

// childName 学生名（容错：不存在返回空）。
func (s *sStudyPlanet) childName(ctx context.Context, id int) string {
	v, err := g.DB().Model("children").Ctx(ctx).Where("id", id).Value("name")
	if err != nil || v == nil {
		return ""
	}
	return v.String()
}

// childAvatar 学生头像。
func (s *sStudyPlanet) childAvatar(ctx context.Context, id int) string {
	v, err := g.DB().Model("children").Ctx(ctx).Where("id", id).Value("avatar")
	if err != nil || v == nil {
		return "🐣"
	}
	return v.String()
}
