package v1

import "github.com/gogf/gf/v2/frame/g"

// BattleRankEntry 段位榜条目。
type BattleRankEntry struct {
	Rank      int    `json:"rank"`
	ChildID   int    `json:"child_id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Trophies  int    `json:"trophies"`
	Tier      string `json:"tier"`
	TierEmoji string `json:"tier_emoji"`
	Wins      int    `json:"wins"`
	Losses    int    `json:"losses"`
	Battles   int    `json:"battles"`
}

// BattleRankReq 对战段位排行榜（按奖杯数）。
type BattleRankReq struct {
	g.Meta    `path:"/battle/rank" method:"get" tags:"Battle" summary:"对战段位榜"`
	Limit     int `json:"limit" in:"query"`
	StudentID int `json:"student_id" in:"query"`
}
type BattleRankRes struct {
	Entries []BattleRankEntry `json:"entries"`
	My      *BattleRankEntry  `json:"my,omitempty"`
}

// BattleHistoryEntry 历史对战条目。
type BattleHistoryEntry struct {
	ID          int    `json:"id"`
	Opponent    string `json:"opponent"`
	OpponentAvatar string `json:"opponent_avatar"`
	MyScore     int    `json:"my_score"`
	OppScore    int    `json:"opp_score"`
	Result      string `json:"result"` // win / lose / draw
	Trophies    int    `json:"trophies"`
	CreatedAt   string `json:"created_at"`
}

// BattleHistoryReq 我的历史对战。
type BattleHistoryReq struct {
	g.Meta    `path:"/battle/history" method:"get" tags:"Battle" summary:"对战历史"`
	StudentID int `json:"student_id" in:"query"`
}
type BattleHistoryRes []BattleHistoryEntry
