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
	Id          interface{} //
	Title       interface{} //
	Type        interface{} //
	DueDate     *gtime.Time //
	Points      interface{} //
	Status      interface{} //
	ChildId     interface{} //
	CreatedAt   *gtime.Time //
	CompletedAt *gtime.Time //
}
