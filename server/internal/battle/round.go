// Package battle 对局过程：推题、倒计时、判分计分、机器人答题。
package battle

import (
	"time"

	"github.com/gogf/gf/v2/util/grand"

	"studyplanet/internal/judge"
)

// nextQuestion 推进到下一题并开始 10 秒计时。
func (e *Engine) nextQuestion(rm *room) {
	rm.mu.Lock()
	if rm.finished {
		rm.mu.Unlock()
		return
	}
	rm.qIndex++ // 前移到下一题（首题从 -1 → 0）
	if rm.qIndex >= len(rm.qs) {
		rm.mu.Unlock()
		e.finishRoom(rm)
		return
	}
	idx := rm.qIndex
	q := rm.qs[idx]
	rm.qDeadli = time.Now().Add(secondsPerQ * time.Second)
	rm.timerOn = true
	rm.mu.Unlock()

	// 倒计时推进（真人客户端各自渲染，服务端 tick 供校准 + 超时判定）
	go e.runTimer(rm, idx)

	for _, p := range []*player{rm.p1, rm.p2} {
		if p.isBot {
			continue
		}
		e.send(p, srvMsg{Type: "question_next", QIndex: idx, Question: q.pub, Remain: secondsPerQ})
	}
	// 机器人答题排程
	if rm.p2.isBot || rm.p1.isBot {
		bot := rm.p2
		if rm.p1.isBot {
			bot = rm.p1
		}
		e.scheduleBot(rm, bot, idx)
	}
}

// runTimer 每秒广播 tick，超时未全员作答则强制进入下一题。
func (e *Engine) runTimer(rm *room, idx int) {
	for {
		time.Sleep(time.Second)
		rm.mu.Lock()
		if rm.finished || rm.qIndex != idx || !rm.timerOn {
			rm.mu.Unlock()
			return
		}
		remain := int(time.Until(rm.qDeadli).Seconds())
		if remain <= 0 {
			rm.timerOn = false
			rm.mu.Unlock()
			e.timeoutQuestion(rm, idx)
			return
		}
		rm.mu.Unlock()
		for _, p := range []*player{rm.p1, rm.p2} {
			if !p.isBot {
				e.send(p, srvMsg{Type: "tick", QIndex: idx, Remain: remain})
			}
		}
	}
}

// scoreFor 剩余秒数 → 得分：满分 10，线性衰减到最低 4 分。
func scoreFor(remainSec int) int {
	s := minScorePerQ + (scorePerQ-minScorePerQ)*remainSec/secondsPerQ
	if s > scorePerQ {
		s = scorePerQ
	}
	if s < minScorePerQ {
		s = minScorePerQ
	}
	return s
}

// onAnswer 处理玩家作答：判分、计分、广播结果；双方都答完或超时 → 下一题。
func (e *Engine) onAnswer(p *player, m cliMsg) {
	// roomOf 内部自带 e.mu 加锁；此处不可先持 e.mu（sync.Mutex 不可重入会死锁）
	rm := e.roomOf(p)
	if rm == nil {
		return
	}
	rm.mu.Lock()
	idx := m.QIndex
	if rm.finished || idx != rm.qIndex || p.answers[idx] || !rm.timerOn {
		rm.mu.Unlock()
		return
	}
	remain := int(time.Until(rm.qDeadli).Seconds())
	if remain < 0 {
		remain = 0
	}
	q := rm.qs[idx]
	correct := judge.Judge(q.answer, m.Answer)
	p.answers[idx] = true
	if correct {
		sc := scoreFor(remain)
		p.scores[idx] = sc
		p.total += sc
		p.gotCount++
		p.streak++
	} else {
		p.streak = 0
	}
	bothDone := rm.p1.answers[idx] && rm.p2.answers[idx]
	opp := otherOf(rm, p)
	oppAnswered := opp.answers[idx]
	oppTotal := opp.total
	rm.mu.Unlock()

	e.send(p, srvMsg{
		Type: "answer_result", QIndex: idx, Correct: correct,
		Score: p.scores[idx], Total: p.total, OppTotal: oppTotal, OppAnswered: oppAnswered,
	})

	if bothDone {
		rm.mu.Lock()
		rm.timerOn = false
		rm.mu.Unlock()
		e.nextQuestion(rm)
	}
}

// timeoutQuestion 单题超时：未作答方计 0 分，进入下一题。
func (e *Engine) timeoutQuestion(rm *room, idx int) {
	for _, p := range []*player{rm.p1, rm.p2} {
		if p.isBot {
			continue
		}
		rm.mu.Lock()
		missed := !p.answers[idx]
		oppTotal := otherOf(rm, p).total
		rm.mu.Unlock()
		if missed {
			p.streak = 0
			e.send(p, srvMsg{Type: "answer_result", QIndex: idx, Correct: false, Score: 0, Total: p.total, OppTotal: oppTotal})
		}
	}
	e.nextQuestion(rm)
}

// scheduleBot 机器人按难度延时作答：大概率答对，速度随机 3~8 秒。
func (e *Engine) scheduleBot(rm *room, bot *player, idx int) {
	delay := time.Duration(3+grand.N(0, 5)) * time.Second
	time.AfterFunc(delay, func() {
		rm.mu.Lock()
		if rm.finished || rm.qIndex != idx || bot.answers[idx] {
			rm.mu.Unlock()
			return
		}
		correct := grand.N(1, 100) <= 62 // 机器人 62% 正确率
		bot.answers[idx] = true
		remain := int(time.Until(rm.qDeadli).Seconds())
		if remain < 0 {
			remain = 0
		}
		if correct {
			sc := scoreFor(remain)
			bot.scores[idx] = sc
			bot.total += sc
			bot.gotCount++
			bot.streak++
		} else {
			bot.streak = 0
		}
		bothDone := rm.p1.answers[idx] && rm.p2.answers[idx]
		rm.mu.Unlock()

		// 给真人对手播报机器人完成状态（只更新 opp_total，不携带 Remain 避免干扰倒计时）
		human := rm.p1
		if rm.p1.isBot {
			human = rm.p2
		}
		e.send(human, srvMsg{Type: "opp_done", QIndex: idx, OppTotal: bot.total})
		if bothDone {
			rm.mu.Lock()
			rm.timerOn = false
			rm.mu.Unlock()
			e.nextQuestion(rm)
		}
	})
}

// roomOf 找玩家所在房间。
func (e *Engine) roomOf(p *player) *room {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, rm := range e.rooms {
		if rm.p1 == p || rm.p2 == p {
			return rm
		}
	}
	return nil
}

// leave 掉线清理：移出等待队列并关闭连接。
// 只关连接（writePump 会随之退出），不 drain channel：
// 原来的 <-p.send 读一行会静默吞掉待发消息，还可能与 close 竞争。
func (e *Engine) leave(p *player) {
	e.mu.Lock()
	for k, v := range e.waiting {
		if v == p {
			delete(e.waiting, k)
		}
	}
	e.mu.Unlock()
	p.conn.Close()
}
