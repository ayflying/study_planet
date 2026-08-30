// Package battle 匹配与开房：队列配对、机器人兜底、房间创建。
package battle

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	v1 "studyplanet/api/studyplanet/v1"
)

// joinQueue 加入匹配队列；同 subject|grade 已有等待者则配对开房，否则入队等 3 秒进机器人局。
func (e *Engine) joinQueue(p *player, m cliMsg) {
	if m.StudentID <= 0 {
		m.StudentID = 1
	}
	if strings.TrimSpace(m.Subject) == "" {
		m.Subject = "math"
	}
	if m.Grade <= 0 {
		m.Grade = 3
	}
	p.childID = m.StudentID
	p.name, p.avatar = e.childInfo(m.StudentID)
	log.Printf("[battle] join student=%d subject=%s grade=%d", p.childID, m.Subject, m.Grade)

	key := fmt.Sprintf("%s|%d", m.Subject, m.Grade)
	e.mu.Lock()
	if other, ok := e.waiting[key]; ok && other != p {
		delete(e.waiting, key)
		e.mu.Unlock()
		e.startRoom(other, p, m.Subject, m.Grade, false)
		return
	}
	e.waiting[key] = p
	e.mu.Unlock()

	// 3 秒内没人配对 → 机器人兜底
	time.AfterFunc(matchWaitBot, func() {
		e.mu.Lock()
		cur, ok := e.waiting[key]
		if ok && cur == p {
			delete(e.waiting, key)
			e.mu.Unlock()
			log.Printf("[battle] bot fallback for student=%d", p.childID)
			bot := &player{
				childID: -1,
				name:    e.names[grand.N(0, len(e.names)-1)] + "（机器人）",
				avatar:  "🤖",
				isBot:   true,
			}
			e.startRoom(p, bot, m.Subject, m.Grade, true)
			return
		}
		if ok {
			log.Printf("[battle] bot fallback skipped: key held by another player")
		}
		e.mu.Unlock()
	})
}

// childInfo 查学生名字与头像（查不到给默认）。
func (e *Engine) childInfo(childID int) (string, string) {
	row, err := g.DB().Model("children").Ctx(gctxNew()).Where("id", childID).One()
	if err != nil || row.IsEmpty() {
		return "小勇士", "🧑‍🚀"
	}
	name := row["name"].String()
	if name == "" {
		name = "小勇士"
	}
	avatar := row["avatar"].String()
	if avatar == "" {
		avatar = "🧑‍🚀"
	}
	return name, avatar
}

// startRoom 开房：抽 5 题、通知双方、启动第一题计时。
// botRoom=true 时 p2 是机器人（自动答题）。
func (e *Engine) startRoom(p1, p2 *player, subject string, grade int, botRoom bool) {
	log.Printf("[battle] startRoom subject=%s grade=%d bot=%v", subject, grade, botRoom)
	rm := &room{
		ID:      fmt.Sprintf("r%d", time.Now().UnixNano()),
		Subject: subject,
		Grade:   grade,
		p1:      p1,
		p2:      p2,
		qIndex:  -1, // nextQuestionFrom 首次调用前移到 0
	}
	qs, err := pickBattleQuestions(subject, grade, questionCount)
	if err != nil {
		log.Printf("[battle] pickBattleQuestions 失败: %v", err)
	}
	if err != nil || len(qs) == 0 {
		e.send(p1, srvMsg{Type: "finished", Result: "error"})
		return
	}
	rm.qs = qs
	e.mu.Lock()
	e.rooms[rm.ID] = rm
	e.mu.Unlock()

	qsForCli := make([]*v1.PubQuestion, len(qs))
	for i, q := range qs {
		qsForCli[i] = q.pub
	}
	for _, p := range []*player{p1, p2} {
		if p.isBot {
			continue
		}
		opp := otherOf(rm, p)
		e.send(p, srvMsg{
			Type:      "matched",
			Room:      rm.ID,
			Opponent:  &oppInfo{Name: opp.name, Avatar: opp.avatar, IsBot: opp.isBot},
			Questions: qsForCli,
		})
	}
	e.nextQuestionFrom(rm, -1)
}

// otherOf 取对方玩家。
func otherOf(rm *room, p *player) *player {
	if rm.p1 == p {
		return rm.p2
	}
	return rm.p1
}

// pickBattleQuestions 从内容库抽题（enabled=1，年级 ±1）。
func pickBattleQuestions(subject string, grade, limit int) ([]*battleQuestion, error) {
	all, err := g.DB().Model("questions").Ctx(gctxNew()).
		Where("subject", subject).Where("enabled", 1).
		WhereGTE("grade", grade-1).WhereLTE("grade", grade+1).
		OrderRandom().Limit(limit).All()
	if err != nil {
		return nil, err
	}
	out := make([]*battleQuestion, 0, len(all))
	for _, r := range all {
		var opts []string
		_ = jsonUnmarshal(r["options"].String(), &opts)
		out = append(out, &battleQuestion{
			pub: &v1.PubQuestion{
				ID:         r["id"].Int(),
				Subject:    r["subject"].String(),
				Grade:      r["grade"].Int(),
				Topic:      r["topic"].String(),
				QType:      r["qtype"].String(),
				Passage:    r["passage"].String(),
				Question:   r["question"].String(),
				Options:    opts,
				Difficulty: r["difficulty"].Int(),
			},
			answer: r["answer"].String(),
		})
	}
	return out, nil
}
