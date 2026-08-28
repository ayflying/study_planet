// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// MathProblems is the golang structure of table math_problems for DAO operations like Where/Data.
type MathProblems struct {
	g.Meta      `orm:"table:math_problems, do:true"`
	Id          any //
	Level       any //
	Type        any //
	Question    any //
	Options     any //
	Answer      any //
	Explanation any //
}
