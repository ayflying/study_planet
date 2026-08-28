// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tasks is the golang structure of table tasks for DAO operations like Where/Data.
type Tasks struct {
	g.Meta      `orm:"table:tasks, do:true"`
	Id          any         //
	Title       any         //
	Type        any         //
	DueDate     *gtime.Time //
	Points      any         //
	Status      any         //
	ChildId     any         //
	CreatedAt   *gtime.Time //
	CompletedAt *gtime.Time //
}
