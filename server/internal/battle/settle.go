// Package battle 结算与持久化：定胜负、奖杯/段位落库、奖励文案。
package battle

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// finishRoom 结算：定胜负、奖杯/段位落库、广播结算画面数据、清理房间。
func (e *Engine) finishRoom(rm *room) {
	rm.mu.Lock()
	if rm.finished {
		rm.mu.Unlock()
		return
	}
	rm.finished = true
	p1, p2 := rm.p1, rm.p2
	rm.mu.Unlock()

	resultP1, resultP2 := "draw", "draw"
	switch {
	case p1.total > p2.total:
		resultP1, resultP2 = "win", "lose"
	case p1.total < p2.total:
		resultP1, resultP2 = "lose", "win"
	}

	// 持久化（机器人不落库）
	ctx := gctxNew()
	var t1, t2 int
	if !p1.isBot {
		t1 = e.persistResult(ctx, rm, p1, resultP1, p2)
	}
	if !p2.isBot {
		t2 = e.persistResult(ctx, rm, p2, resultP2, p1)
	}

	for _, ps := range []*player{p1, p2} {
		if ps.isBot {
			continue
		}
		res, trophies := resultP1, t1
		if ps == p2 {
			res, trophies = resultP2, t2
		}
		tier, tierEmoji := TierName(trophies)
		e.send(ps, srvMsg{
			Type: "finished", Result: res,
			MyScore: ps.total, OppScore: otherOf(rm, ps).total,
			Trophies: trophies, Tier: tier, TierEmoji: tierEmoji,
			WinStreak: ps.streak,
			Exp:       expReward(res),
			Rewards:   rewardCopy(res, trophies, ps.streak),
			Snack:     e.snackOf(res, ps.childID),
		})
	}

	e.mu.Lock()
	delete(e.rooms, rm.ID)
	e.mu.Unlock()
	// 关闭连接，客户端拿到 finished 后自行跳结算页
	for _, ps := range []*player{p1, p2} {
		closeSend(ps)
	}
	log.Printf("[battle] room %s finished p1=%d p2=%d", rm.ID, p1.total, p2.total)
}

// snackOf 胜利时触发零食掉落，返回掉落零食 id（空=未掉落）。
func (e *Engine) snackOf(result string, childID int) string {
	if result == "win" && e.OnSnack != nil {
		return e.OnSnack(childID)
	}
	return ""
}

// closeSend 安全关闭玩家发送通道（只关一次，防 close/send 竞争）。
func closeSend(p *player) {
	if p.isBot {
		return
	}
	p.closed.Lock()
	if !p.closedFlag {
		p.closedFlag = true
		close(p.send)
	}
	p.closed.Unlock()
}

// expReward 结算经验：胜 30 / 平 15 / 负 8。
func expReward(result string) int {
	switch result {
	case "win":
		return 30
	case "draw":
		return 15
	default:
		return 8
	}
}

// rewardCopy 结算奖励文案（多档位+连击加成+段位里程碑）。
func rewardCopy(result string, trophies, streak int) []string {
	var rs []string
	switch result {
	case "win":
		rs = append(rs, "🏆 对战胜利 +20 奖杯", "⭐ 经验 +30", "🎁 战胜宝箱 ×1（家长奖励池随机掉落）")
		if streak >= 2 {
			rs = append(rs, fmt.Sprintf("🔥 %d 连击！额外经验 +10", streak))
		}
	case "draw":
		rs = append(rs, "🤝 势均力敌 +5 奖杯", "⭐ 经验 +15", "📦 安慰礼包 ×1")
	default:
		rs = append(rs, "💪 虽败犹荣（奖杯不清零）", "⭐ 经验 +8", "🍀 复仇之火：下一场胜利奖杯 ×2")
	}
	switch {
	case trophies >= 300:
		rs = append(rs, "👑 已达王者段位，全服仰望！")
	case trophies >= 220:
		rs = append(rs, "🌟 距离王者仅一步之遥！")
	case trophies >= 150:
		rs = append(rs, "💎 钻石段位达成！")
	case trophies >= 100:
		rs = append(rs, "🛡️ 铂金段位达成！")
	case trophies >= 60:
		rs = append(rs, "🏅 黄金段位达成！")
	case trophies >= 30:
		rs = append(rs, "🥈 白银段位达成！")
	}
	return rs
}

// persistResult 战绩落库：battle_scores 奖杯/胜场更新 + battles 记录一场对局。
func (e *Engine) persistResult(ctx context.Context, rm *room, p *player, result string, opp *player) int {
	delta := TrophiesDelta(result)

	// battle_scores：先查后写（兼容多数据库）
	cur, err := g.DB().Model("battle_scores").Ctx(ctx).Where("child_id", p.childID).One()
	if err != nil {
		log.Printf("battle_scores 查询失败 child=%d: %v", p.childID, err)
		return 0
	}
	trophies := 0
	if cur.IsEmpty() {
		trophies = maxInt(0, delta)
		wins, losses, draws := 0, 1, 0
		switch result {
		case "win":
			wins, losses = 1, 0
		case "draw":
			wins, draws, losses = 0, 1, 0
		}
		streak := 0
		if result == "win" {
			streak = 1
		}
		_, _ = g.DB().Model("battle_scores").Ctx(ctx).Data(g.Map{
			"child_id":    p.childID,
			"trophies":    trophies,
			"wins":        wins,
			"losses":      losses,
			"draws":       draws,
			"battles":     1,
			"best_streak": streak,
			"cur_streak":  streak,
		}).Insert()
	} else {
		trophies = maxInt(0, cur["trophies"].Int()+delta)
		cs := cur["cur_streak"].Int()
		if result != "win" {
			cs = 0
		} else {
			cs++
		}
		updates := g.Map{
			"trophies":   trophies,
			"battles":    cur["battles"].Int() + 1,
			"cur_streak": cs,
		}
		switch result {
		case "win":
			updates["wins"] = cur["wins"].Int() + 1
			if cs > cur["best_streak"].Int() {
				updates["best_streak"] = cs
			}
		case "draw":
			updates["draws"] = cur["draws"].Int() + 1
		default:
			updates["losses"] = cur["losses"].Int() + 1
		}
		_, _ = g.DB().Model("battle_scores").Ctx(ctx).Where("child_id", p.childID).Data(updates).Update()
	}

	// battles 对局记录
	qids := make([]string, 0, len(rm.qs))
	for _, q := range rm.qs {
		qids = append(qids, fmt.Sprintf("%d", q.pub.ID))
	}
	oppID := opp.childID
	_, _ = g.DB().Model("battles").Ctx(ctx).Data(g.Map{
		"room_id":      rm.ID,
		"p1_id":        p.childID, // 简化：每方各记一条，p1=本人 p2=对手
		"p2_id":        oppID,
		"subject":      rm.Subject,
		"grade":        rm.Grade,
		"status":       "finished",
		"winner_id":    map[string]int{"win": p.childID, "draw": 0, "lose": oppID}[result],
		"p1_score":     p.total,
		"p2_score":     opp.total,
		"p1_correct":   p.gotCount,
		"p2_correct":   opp.gotCount,
		"p1_robot":     boolInt(p.isBot),
		"p2_robot":     boolInt(opp.isBot),
		"trophies1":    delta,
		"question_ids": strings.Join(qids, ","),
		"finished_at":  gtime.Now(),
	}).Insert()

	// 经验奖励（联动周榜）
	if e.AddXP != nil {
		e.AddXP(p.childID, expReward(result))
	}
	return trophies
}
