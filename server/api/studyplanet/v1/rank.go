package v1

import "github.com/gogf/gf/v2/frame/g"

// LeaderboardEntry 周榜条目。
type LeaderboardEntry struct {
	Rank    int    `json:"rank"`
	ChildID int    `json:"child_id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	XP      int    `json:"xp"`
}

// WeeklyLeaderboardReq 每周经验排行榜（当前 ISO 周）。
type WeeklyLeaderboardReq struct {
	g.Meta `path:"/leaderboard/weekly" method:"get" tags:"Leaderboard" summary:"周经验排行榜"`
	Limit  int `json:"limit" in:"query"`
}
type WeeklyLeaderboardRes struct {
	Week    string             `json:"week"`
	Redis   bool               `json:"redis"`
	Entries []LeaderboardEntry `json:"entries"`
	MyXP    int                `json:"my_xp"`
	MyRank  int                `json:"my_rank"`
	MyID    int                `json:"my_id"`
}

// WrongQuestion 错题本条目。
type WrongQuestion struct {
	ID              int    `json:"id"`
	ChildID         int    `json:"child_id"`
	Subject         string `json:"subject"`
	RefID           int    `json:"ref_id"`
	WrongCount      int    `json:"wrong_count"`
	Resolved        int    `json:"resolved"`
	LastWrongAt     string `json:"last_wrong_at"`
	LastReviewedAt  string `json:"last_reviewed_at"`
}

// ListWrongQuestionsReq 学生错题本（可选按 subject 过滤）。
type ListWrongQuestionsReq struct {
	g.Meta  `path:"/wrong-questions" method:"get" tags:"WrongBook" summary:"错题本"`
	Subject string `json:"subject" in:"query"`
}
type ListWrongQuestionsRes []WrongQuestion
