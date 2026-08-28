// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LeaderboardWeekly is the golang structure for table leaderboard_weekly.
type LeaderboardWeekly struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	WeekKey   string      `json:"weekKey"   orm:"week_key"   description:""` //
	ChildId   int64       `json:"childId"   orm:"child_id"   description:""` //
	Xp        int         `json:"xp"        orm:"xp"         description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
