// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PracticeSessions is the golang structure of table practice_sessions for DAO operations like Where/Data.
type PracticeSessions struct {
	g.Meta     `orm:"table:practice_sessions, do:true"`
	Id         any         //
	ChildId    any         //
	Subject    any         //
	Level      any         //
	Total      any         //
	Correct    any         //
	MaxCombo   any         //
	Bonus      any         //
	Stars      any         //
	Finished   any         //
	CreatedAt  *gtime.Time //
	FinishedAt *gtime.Time //
}
