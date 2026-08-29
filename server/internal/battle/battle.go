// Package battle 真人对战引擎：WebSocket 匹配 → 房间 → 5题×10秒速度计分 → 结算段位。
//
// 玩法规则（产品需求）：
//   - 每场 5 题，每题 10 秒，无应用题（题库层面已保证）；
//   - 每题 10 分，答对越快加分越多（按剩余时间比例给 4~10 分）；
//   - 分高者胜，胜 +20 奖杯、平 +5、负 -10（下限 0）；
//   - 匹配不到真人时 3 秒后机器人兜底，保证随时可玩。
//
// 协议（JSON over WS）：
//   客户端 → 服务端：{type:"join", student_id, subject, grade}
//                   {type:"answer", qindex, answer}
//   服务端 → 客户端：{type:"matched", room, opponent:{name,avatar}, questions:[...]}
//                   {type:"tick", qindex, remain}
//                   {type:"answer_result", qindex, correct, score, total, opp_total}
//                   {type:"question_next", qindex, question}
//                   {type:"finished", result, my_score, opp_score, trophies, tier...}
package battle

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gorilla/websocket"

	v1 "studyplanet/api/studyplanet/v1"
	"studyplanet/internal/judge"
)

// ---------- 常量 ----------

const (
	questionCount   = 5               // 每场题数
	secondsPerQ     = 10              // 每题答题时长（秒）
	scorePerQ       = 10              // 每题满分
	minScorePerQ    = 4               // 答对最低得分（掐点答对）
	matchWaitBot    = 3 * time.Second // 真人匹配等待时长，超时进机器人
	perQuestionTick = 250 * time.Millisecond
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = 50 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true }, // 同源由前端部署保证，教学项目放开
}

// ---------- 消息结构 ----------

type cliMsg struct {
	Type      string `json:"type"`
	StudentID int    `json:"student_id"`
	Subject   string `json:"subject"`
	Grade     int    `json:"grade"`
	QIndex    int    `json:"qindex"`
	Answer    string `json:"answer"`
}

type oppInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	IsBot  bool   `json:"is_bot"`
}

type srvMsg struct {
	Type        string           `json:"type"`
	Room        string           `json:"room,omitempty"`
	Opponent    *oppInfo         `json:"opponent,omitempty"`
	Questions   []*v1.PubQuestion `json:"questions,omitempty"`
	Question    *v1.PubQuestion  `json:"question,omitempty"`
	QIndex      int              `json:"qindex,omitempty"`
	Remain      int              `json:"remain,omitempty"`
	Correct     bool             `json:"correct,omitempty"`
	Score       int              `json:"score,omitempty"`
	Total       int              `json:"total,omitempty"`
	OppTotal    int              `json:"opp_total,omitempty"`
	OppAnswered bool             `json:"opp_answered,omitempty"`
	Result      string           `json:"result,omitempty"` // win/lose/draw
	MyScore     int              `json:"my_score,omitempty"`
	OppScore    int              `json:"opp_score,omitempty"`
	Trophies    int              `json:"trophies,omitempty"`
	Tier        string           `json:"tier,omitempty"`
	TierEmoji   string           `json:"tier_emoji,omitempty"`
	WinStreak   int              `json:"win_streak,omitempty"`
	Rewards     []string         `json:"rewards,omitempty"` // 结算奖励文案列表
	Exp         int              `json:"exp,omitempty"`     // 结算经验
}

// ---------- 玩家与房间 ----------

type player struct {
	childID int
	name    string
	avatar  string
	conn    *websocket.Conn
	send    chan []byte
	isBot   bool

	closed     sync.Mutex // 保护 closedFlag（close 与 send 竞争防护）
	closedFlag bool

	// 对战运行时
	answers  [questionCount]bool // 每题是否已作答
	scores   [questionCount]int  // 每题得分
	total    int
	gotCount int
	streak   int // 本场连对（结算奖励用）
}

type room struct {
	ID       string
	Subject  string
	Grade    int
	p1, p2   *player
	qs       []*battleQuestion // 题目（含答案，服务端持有）
	qIndex   int               // 当前题下标（nextQuestion 前移，起始 -1 表示尚未开题）
	qDeadli  time.Time
	timerOn  bool
	finished bool
	mu       sync.Mutex
}

type battleQuestion struct {
	pub    *v1.PubQuestion
	answer string
}

// ---------- 引擎 ----------

// Engine 对战引擎：管理匹配队列与房间生命周期。
type Engine struct {
	mu       sync.Mutex
	waiting  map[string]*player // key: subject|grade 简单分级匹配
	rooms    map[string]*room
	names    []string // 机器人名字池
	// AddXP 经验回调：由 cmd 启动时注入（联动周榜），可空。
	AddXP func(childID int, delta int)
}

// New 创建引擎。
func New() *Engine {
	return &Engine{
		waiting: map[string]*player{},
		rooms:   map[string]*room{},
		names:   []string{"小飞侠", "闪电麦昆", "学霸猫", "题海冲浪手", "口算小王子", "智慧星", "无敌小钢炮", "风一样的少年", "数字猎人", "小神算"},
	}
}

// HandleWS HTTP → WebSocket 升级入口（挂到 ghttp）。
func (e *Engine) HandleWS(r *ghttp.Request) {
	conn, err := upgrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		return
	}
	p := &player{conn: conn, send: make(chan []byte, 32)}
	go e.writePump(p)
	// join 消息驱动后续流程
	for {
		var m cliMsg
		if err := conn.ReadJSON(&m); err != nil {
			e.leave(p)
			return
		}
		switch m.Type {
		case "join":
			e.joinQueue(p, m)
		case "answer":
			e.onAnswer(p, m)
		}
	}
}

// writePump 单独 goroutine 写 WS，避免并发写冲突。
func (e *Engine) writePump(p *player) {
	defer func() {
		p.conn.Close()
	}()
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case msg, ok := <-p.send:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := p.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// send 向玩家投递消息（非阻塞）；已关闭连接安全跳过，绝不 panic。
func (e *Engine) send(p *player, m srvMsg) {
	if p.isBot {
		return
	}
	b, err := jsonMarshal(m)
	if err != nil {
		return
	}
	select {
	case p.send <- b:
	default: // 队列满视为掉线
		e.leave(p)
	}
}

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
		qIndex:  -1, // nextQuestion 首次调用前移到 0
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
		e.send(p, srvMsg{
			Type:      "matched",
			Room:      rm.ID,
			Opponent:  &oppInfo{Name: otherOf(rm, p).name, Avatar: otherOf(rm, p).avatar, IsBot: otherOf(rm, p).isBot},
			Questions: qsForCli,
		})
	}
	e.nextQuestion(rm)
}

// otherOf 取对方玩家。
func otherOf(rm *room, p *player) *player {
	if rm.p1 == p {
		return rm.p2
	}
	return rm.p1
}

// pickBattleQuestions 从内容库抽题（enabled=1，年级 ±1，优先填空+选择混合）。
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
		q := rm.qs[idx]
		correct := grand.N(1, 100) <= 62 // 机器人 62% 正确率
		_ = q // 答案仅服务端持有，机器人无需真实作答内容
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

	var resultP1, resultP2 string
	switch {
	case p1.total > p2.total:
		resultP1, resultP2 = "win", "lose"
	case p1.total < p2.total:
		resultP1, resultP2 = "lose", "win"
	default:
		resultP1, resultP2 = "draw", "draw"
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
	_ = t2

	for _, ps := range []*player{p1, p2} {
		if ps.isBot {
			continue
		}
		res := resultP1
		trophies := t1
		if ps == p2 {
			res = resultP2
			trophies = t2
		}
		tier, tierEmoji := tierName(trophies)
		e.send(ps, srvMsg{
			Type: "finished", Result: res,
			MyScore: ps.total, OppScore: otherOf(rm, ps).total,
			Trophies: trophies, Tier: tier, TierEmoji: tierEmoji,
			WinStreak: ps.streak,
			Exp:       expReward(res),
			Rewards:   rewardCopy(res, trophies, ps.streak),
		})
	}

	e.mu.Lock()
	delete(e.rooms, rm.ID)
	e.mu.Unlock()
	// 关闭连接，客户端拿到 finished 后自行跳结算页
	for _, ps := range []*player{p1, p2} {
		if !ps.isBot {
			ps.closed.Lock()
			if !ps.closedFlag {
				ps.closedFlag = true
				close(ps.send)
			}
			ps.closed.Unlock()
		}
	}
	log.Printf("[battle] room %s finished p1=%d p2=%d", rm.ID, p1.total, p2.total)
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

// rewardCopy 结算奖励文案（奖励好一点：多档位+连击加成+段位里程碑）。
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

// tierName 奖杯 → 段位（与 battle_api.go 同规则，这里独立实现避免循环依赖）。
func tierName(trophies int) (string, string) {
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

// persistResult 战绩落库：battle_scores 奖杯/胜场更新 + battles 记录一场对局。
func (e *Engine) persistResult(ctx context.Context, rm *room, p *player, result string, opp *player) int {
	delta := map[string]int{"win": 20, "draw": 5, "lose": -10}[result]

	// battle_scores upsert（SQLite 与 MySQL 兼容写法：先查后写）
	cur, err := g.DB().Model("battle_scores").Ctx(ctx).Where("child_id", p.childID).One()
	if err != nil {
		log.Printf("battle_scores 查询失败 child=%d: %v", p.childID, err)
		return 0
	}
	trophies := 0
	if cur.IsEmpty() {
		trophies = maxInt(0, delta)
		wins, losses, draws, battles := 0, 0, 0, 1
		switch result {
		case "win":
			wins = 1
		case "draw":
			draws = 1
		default:
			losses = 1
		}
		streak := 0
		if result == "win" {
			streak = 1
		}
		_, _ = g.DB().Model("battle_scores").Ctx(ctx).Data(g.Map{
			"child_id":   p.childID,
			"trophies":   trophies,
			"wins":       wins,
			"losses":     losses,
			"draws":      draws,
			"battles":    battles,
			"best_streak": streak,
			"cur_streak": streak,
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
		"room_id":    rm.ID,
		"p1_id":      p.childID, // 简化：每方各记一条，p1=本人 p2=对手
		"p2_id":      oppID,
		"subject":    rm.Subject,
		"grade":      rm.Grade,
		"status":     "finished",
		"winner_id":  map[string]int{"win": p.childID, "draw": 0, "lose": oppID}[result],
		"p1_score":   p.total,
		"p2_score":   opp.total,
		"p1_correct": p.gotCount,
		"p2_correct": opp.gotCount,
		"p1_robot":   boolInt(p.isBot),
		"p2_robot":   boolInt(opp.isBot),
		"trophies1":  delta,
		"question_ids": strings.Join(qids, ","),
		"finished_at": gtime.Now(),
	}).Insert()

	// 经验奖励（联动周榜）
	if e.AddXP != nil {
		e.AddXP(p.childID, expReward(result))
	}
	return trophies
}

// ---------- 掉线清理 ----------

func (e *Engine) leave(p *player) {
	e.mu.Lock()
	for k, v := range e.waiting {
		if v == p {
			delete(e.waiting, k)
		}
	}
	e.mu.Unlock()
	// 只关连接（writePump 会随之退出），不 drain channel：
	// 原来的 <-p.send 读一行会静默吞掉待发消息，还可能与 close 竞争。
	p.conn.Close()
}

// ---------- 辅助 ----------

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
