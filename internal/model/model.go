package model

// 数据模型：字段映射数据库列（db tag），JSON 用于 API 输出。

type Word struct {
	ID        int    `db:"id" json:"id"`
	Level     int    `db:"level" json:"level"`
	Word      string `db:"word" json:"word"`
	Meaning   string `db:"meaning" json:"meaning"`
	Phonetic  string `db:"phonetic" json:"phonetic"`
	Example   string `db:"example" json:"example"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type WordProgress struct {
	WordID       int    `db:"word_id" json:"word_id"`
	ChildID      int    `db:"child_id" json:"child_id"`
	Known        int    `db:"known" json:"known"`
	LastReviewed string `db:"last_reviewed" json:"last_reviewed"`
}

type Reading struct {
	ID      int    `db:"id" json:"id"`
	Title   string `db:"title" json:"title"`
	Content string `db:"content" json:"content"`
	Level   int    `db:"level" json:"level"`
}

type ReadingQuestion struct {
	ID        int    `db:"id" json:"id"`
	ReadingID int    `db:"reading_id" json:"reading_id"`
	Question  string `db:"question" json:"question"`
	OptionA   string `db:"option_a" json:"option_a"`
	OptionB   string `db:"option_b" json:"option_b"`
	OptionC   string `db:"option_c" json:"option_c"`
	OptionD   string `db:"option_d" json:"option_d"`
	Answer    string `db:"answer" json:"answer"`
}

type MathProblem struct {
	ID          int    `db:"id" json:"id"`
	Level       int    `db:"level" json:"level"`
	Type        string `db:"type" json:"type"`
	Question    string `db:"question" json:"question"`
	Options     string `db:"options" json:"options"` // JSON 数组字符串
	Answer      string `db:"answer" json:"answer"`
	Explanation string `db:"explanation" json:"explanation"`
}

type Task struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Type        string `db:"type" json:"type"`
	DueDate     string `db:"due_date" json:"due_date"`
	Points      int    `db:"points" json:"points"`
	Status      string `db:"status" json:"status"` // pending | done | overdue(计算字段)
	CreatedAt   string `db:"created_at" json:"created_at"`
	CompletedAt string `db:"completed_at" json:"completed_at"`
}

type PointsLog struct {
	ID        int    `db:"id" json:"id"`
	ChildID   int    `db:"child_id" json:"child_id"`
	Delta     int    `db:"delta" json:"delta"`
	Reason    string `db:"reason" json:"reason"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type Reward struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	CostPoints int    `db:"cost_points" json:"cost_points"`
	Status     string `db:"status" json:"status"` // active | redeemed
}

type Redemption struct {
	ID          int    `db:"id" json:"id"`
	RewardID    int    `db:"reward_id" json:"reward_id"`
	ChildID     int    `db:"child_id" json:"child_id"`
	Status      string `db:"status" json:"status"` // pending | confirmed
	RequestedAt string `db:"requested_at" json:"requested_at"`
	ConfirmedAt string `db:"confirmed_at" json:"confirmed_at"`
}

type Setting struct {
	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`
}

// Student 学生档案（children 表；username 可为空，非空则唯一）。
type Student struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Username  string `db:"username" json:"username"`
	Avatar    string `db:"avatar" json:"avatar"`
	Grade     int    `db:"grade" json:"grade"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

// Parent Casdoor SSO 登录后落库的家长账号。
type Parent struct {
	ID          int    `db:"id" json:"id"`
	CasdoorSub  string `db:"casdoor_sub" json:"-"`
	DisplayName string `db:"display_name" json:"display_name"`
	Avatar      string `db:"avatar" json:"avatar"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	LastLoginAt string `db:"last_login_at" json:"last_login_at"`
}

// PracticeSession 多邻国式练习场次（一次关卡）。
type PracticeSession struct {
	ID         int    `db:"id" json:"id"`
	ChildID    int    `db:"child_id" json:"child_id"`
	Subject    string `db:"subject" json:"subject"`
	Level      int    `db:"level" json:"level"`
	Total      int    `db:"total" json:"total"`
	Correct    int    `db:"correct" json:"correct"`
	MaxCombo   int    `db:"max_combo" json:"max_combo"`
	Bonus      int    `db:"bonus" json:"bonus"`
	Stars      int    `db:"stars" json:"stars"`
	Finished   int    `db:"finished" json:"finished"`
	CreatedAt  string `db:"created_at" json:"created_at"`
	FinishedAt string `db:"finished_at" json:"finished_at"`
}
