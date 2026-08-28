package v1

import "github.com/gogf/gf/v2/frame/g"

// PracticeSession 练习场次（多邻国式关卡）。
type PracticeSession struct {
	ID         int    `json:"id"`
	ChildID    int    `json:"child_id"`
	Subject    string `json:"subject"`
	Level      int    `json:"level"`
	Total      int    `json:"total"`
	Correct    int    `json:"correct"`
	MaxCombo   int    `json:"max_combo"`
	Bonus      int    `json:"bonus"`
	Stars      int    `json:"stars"`
	Finished   int    `json:"finished"`
	CreatedAt  string `json:"created_at"`
	FinishedAt string `json:"finished_at"`
}

// CreateSessionReq 开启一关。
type CreateSessionReq struct {
	g.Meta  `path:"/sessions" method:"post" tags:"Practice" summary:"开始练习场次"`
	Subject string `json:"subject" v:"required"`
	Level   int    `json:"level"`
	Total   int    `json:"total"`
}
type CreateSessionRes PracticeSession

// ListSessionsReq 学生最近练习记录（可选按 level/subject 过滤）。
type ListSessionsReq struct {
	g.Meta  `path:"/sessions" method:"get" tags:"Practice" summary:"练习记录"`
	Level   string `json:"level" in:"query"`
	Subject string `json:"subject" in:"query"`
}
type ListSessionsRes []PracticeSession

// FinishSessionReq 结算一关（同一关只结算一次；按正确率给星级与奖励）。
type FinishSessionReq struct {
	g.Meta `path:"/sessions/:id/finish" method:"post" tags:"Practice" summary:"结算练习"`
	ID     int `in:"path" json:"-"`
}
type FinishSessionRes struct {
	Stars     int `json:"stars"`
	Bonus     int `json:"bonus"`
	MaxCombo  int `json:"max_combo"`
	Correct   int `json:"correct,omitempty"`
	Total     int `json:"total,omitempty"`
	XPGained  int `json:"xp_gained,omitempty"`
	Already   int `json:"already,omitempty"` // 重复结算时为 1
}
