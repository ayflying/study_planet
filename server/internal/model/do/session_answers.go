// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SessionAnswers is the golang structure of table session_answers for DAO operations like Where/Data.
type SessionAnswers struct {
	g.Meta     `orm:"table:session_answers, do:true"`
	Id         any         //
	SessionId  any         //
	RefId      any         //
	Correct    any         //
	Points     any         //
	Combo      any         //
	AnsweredAt *gtime.Time //
}
