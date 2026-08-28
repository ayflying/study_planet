// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure of table questions for DAO operations like Where/Data.
type Questions struct {
	g.Meta      `orm:"table:questions, do:true"`
	Id          any         //
	Subject     any         //
	Grade       any         //
	Topic       any         //
	Qtype       any         //
	Passage     any         //
	Question    any         //
	Options     any         //
	Answer      any         //
	Explanation any         //
	Difficulty  any         //
	Source      any         //
	ContentHash any         //
	Enabled     any         //
	CreatedAt   *gtime.Time //
}
