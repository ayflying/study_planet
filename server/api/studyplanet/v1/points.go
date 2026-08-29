package v1

import "github.com/gogf/gf/v2/frame/g"

// PointsSummaryReq 积分汇总（公开接口）。
type PointsSummaryReq struct {
	g.Meta    `path:"/points" method:"get" tags:"Points" summary:"积分汇总"`
	StudentID int `json:"student_id" in:"query"`
}
type PointsSummaryRes struct {
	Total       int `json:"total"`
	TodayEarned int `json:"today_earned"`
	StudentID   int `json:"student_id"`
}

// PointsLogItem 积分流水条目。
type PointsLogItem struct {
	ID        int    `json:"id"`
	ChildID   int    `json:"child_id"`
	Delta     int    `json:"delta"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

// PointsLogReq 积分流水（最近 100 条，公开接口）。
type PointsLogReq struct {
	g.Meta    `path:"/points/log" method:"get" tags:"Points" summary:"积分流水"`
	StudentID int `json:"student_id" in:"query"`
}
type PointsLogRes []PointsLogItem
