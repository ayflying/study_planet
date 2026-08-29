package studyplanet

import (
	"context"
	"math"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

// tierOf 奖杯数 → 段位（名称 + emoji）。
// 青铜 0+ / 白银 30+ / 黄金 60+ / 铂金 100+ / 钻石 150+ / 星耀 220+ / 王者 300+。
func tierOf(trophies int) (string, string) {
	switch {
	case trophies >= 300:
		return "王者", "👑"
	case trophies >= 220:
		return "星耀", "🌟"
	case trophies >= 150:
		return "钻石", "💎"
	case trophies >= 100:
		return "铂金", "🛡️"
	case trophies >= 60:
		return "黄金", "🏅"
	case trophies >= 30:
		return "白银", "🥈"
	default:
		return "青铜", "🥉"
	}
}

// trophiesDelta 对战奖杯增减：胜 +20、平 +5、负 -10（下限 0）。
func trophiesDelta(result string) int {
	switch result {
	case "win":
		return 20
	case "draw":
		return 5
	default:
		return -10
	}
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
		if r["p2_robot"].Int() == 0 && r["p1_robot"].Int() == 0 && oppID > 0 {
			opName = s.childName(ctx, oppID)
			opAvatar = s.childAvatar(ctx, oppID)
		} else if isP1 && r["p2_robot"].Int() == 0 && oppID > 0 {
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

// applyBattleResult 对战结算落库：更新 battle_scores（奖杯/胜负/连胜）。
// 由 battle 引擎在对局结束时调用（每方各一次）。
func (s *sStudyPlanet) applyBattleResult(ctx context.Context, childID int, result string) {
	if childID <= 0 {
		return
	}
	row, err := g.DB().Model("battle_scores").Ctx(ctx).Where("child_id", childID).One()
	if err != nil {
		gLog("查询 battle_scores 失败: %v", err)
		return
	}
	if row.IsEmpty() {
		if _, err := g.DB().Model("battle_scores").Ctx(ctx).Data(g.Map{"child_id": childID}).Insert(); err != nil {
			gLog("创建 battle_scores 失败: %v", err)
			return
		}
	}
	trophies := row["trophies"].Int()
	wins, losses, draws := row["wins"].Int(), row["losses"].Int(), row["draws"].Int()
	curStreak, bestStreak := row["cur_streak"].Int(), row["best_streak"].Int()
	switch result {
	case "win":
		wins++
		curStreak++
		if curStreak > bestStreak {
			bestStreak = curStreak
		}
	case "draw":
		draws++
	default:
		losses++
		curStreak = 0
	}
	nt := trophies + trophiesDelta(result)
	if nt < 0 {
		nt = 0
	}
	if _, err := g.DB().Model("battle_scores").Ctx(ctx).Where("child_id", childID).Data(g.Map{
		"trophies": nt, "wins": wins, "losses": losses, "draws": draws,
		"cur_streak": curStreak, "best_streak": bestStreak,
		"battles":    row["battles"].Int() + 1,
		"updated_at": gtime.Now(),
	}).Update(); err != nil {
		gLog("更新 battle_scores 失败: %v", err)
	}
}

// battleTrophies 当前奖杯数（无记录为 0）。
func (s *sStudyPlanet) battleTrophies(ctx context.Context, childID int) int {
	v, err := g.DB().Model("battle_scores").Ctx(ctx).Where("child_id", childID).Value("trophies")
	if err != nil || v == nil {
		return 0
	}
	return v.Int()
}

// battleRewardTitle 结算称号（奖励画面用，越赢越稀有）。
func battleRewardTitle(result string, streak, score int) string {
	switch result {
	case "win":
		if score == 50 {
			return "满分传说！⭐"
		}
		if streak >= 3 {
			return "三连胜王者！🔥"
		}
		return "胜利医师！🎉"
	case "draw":
		return "势均力敌！🤝"
	default:
		if score >= 30 {
			return "虽败犹荣！💪"
		}
		return "下次一定能赢！🌱"
	}
}

// battleMaxRewards 奖励上限保护（数值展示用）。
func battleMaxRewards(a, b int) int {
	if a > b {
		return a
	}
	if b > math.MaxInt32 {
		return math.MaxInt32
	}
	return b
}
