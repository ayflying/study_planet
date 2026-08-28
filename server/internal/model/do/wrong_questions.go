// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// WrongQuestions is the golang structure of table wrong_questions for DAO operations like Where/Data.
type WrongQuestions struct {
	g.Meta         `orm:"table:wrong_questions, do:true"`
	Id             any         //
	ChildId        any         //
	Subject        any         //
	RefId          any         //
	WrongCount     any         //
	Resolved       any         //
	LastWrongAt    *gtime.Time //
	LastReviewedAt *gtime.Time //
	CreatedAt      *gtime.Time //
}
