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
//
// 文件分工：protocol.go 结构与常量 / match.go 匹配开房 / round.go 对局过程 /
// settle.go 结算落库 / tier.go 段位规则。
package battle

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Engine 对战引擎：管理匹配队列与房间生命周期。
type Engine struct {
	mu      sync.Mutex
	waiting map[string]*player // key: subject|grade 简单分级匹配
	rooms   map[string]*room
	names   []string // 机器人名字池
	// AddXP 经验回调：由 cmd 启动时注入（联动周榜），可空。
	AddXP func(childID int, delta int)
	// OnSnack 对战胜利零食掉落回调：返回掉落零食 id（空=未掉落）。由 cmd 注入，可空。
	OnSnack func(childID int) string
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
				p.conn.WriteMessage(WebSocketClose, []byte{})
				return
			}
			if err := p.conn.WriteMessage(WebSocketText, msg); err != nil {
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(WebSocketPing, nil); err != nil {
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
