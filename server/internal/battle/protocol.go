// Package battle 协议与数据结构：WS 消息、常量、玩家/房间模型。
package battle

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	v1 "studyplanet/api/studyplanet/v1"
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

// WebSocket 消息类型别名（writePump 使用）。
const (
	WebSocketText  = websocket.TextMessage
	WebSocketPing  = websocket.PingMessage
	WebSocketClose = websocket.CloseMessage
)

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
	Type        string            `json:"type"`
	Room        string            `json:"room,omitempty"`
	Opponent    *oppInfo          `json:"opponent,omitempty"`
	Questions   []*v1.PubQuestion `json:"questions,omitempty"`
	Question    *v1.PubQuestion   `json:"question,omitempty"`
	// 注意：数值/布尔字段一律不加 omitempty——qindex=0、score=0、correct=false
	// 等零值是合法协议值，omitempty 会让字段从 JSON 消失，客户端状态错位。
	QIndex      int               `json:"qindex"`
	Remain      int               `json:"remain,omitempty"`
	Correct     bool              `json:"correct"`
	Score       int               `json:"score"`
	Total       int               `json:"total"`
	OppTotal    int               `json:"opp_total,omitempty"`
	OppAnswered bool              `json:"opp_answered,omitempty"`
	Result      string            `json:"result,omitempty"` // win/lose/draw
	MyScore     int               `json:"my_score"`
	OppScore    int               `json:"opp_score"`
	Trophies    int               `json:"trophies,omitempty"`
	Tier        string            `json:"tier,omitempty"`
	TierEmoji   string            `json:"tier_emoji,omitempty"`
	WinStreak   int               `json:"win_streak,omitempty"`
	Rewards     []string          `json:"rewards,omitempty"` // 结算奖励文案列表
	Exp         int               `json:"exp,omitempty"`     // 结算经验
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
	qIndex   int               // 当前题下标（nextQuestionFrom 前移，起始 -1 表示尚未开题）
	qDeadli  time.Time
	timerOn  bool
	finished bool
	mu       sync.Mutex
}

type battleQuestion struct {
	pub    *v1.PubQuestion
	answer string
}

// ---------- 小工具 ----------

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
