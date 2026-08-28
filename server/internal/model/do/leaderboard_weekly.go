// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LeaderboardWeekly is the golang structure of table leaderboard_weekly for DAO operations like Where/Data.
type LeaderboardWeekly struct {
	g.Meta    `orm:"table:leaderboard_weekly, do:true"`
	Id        any         //
	WeekKey   any         //
	ChildId   any         //
	Xp        any         //
	UpdatedAt *gtime.Time //
}
