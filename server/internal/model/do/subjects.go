// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Subjects is the golang structure of table subjects for DAO operations like Where/Data.
type Subjects struct {
	g.Meta    `orm:"table:subjects, do:true"`
	Id        any         //
	Code      any         //
	Name      any         //
	Icon      any         //
	Color     any         //
	MinGrade  any         //
	MaxGrade  any         //
	Sort      any         //
	Enabled   any         //
	CreatedAt *gtime.Time //
}
