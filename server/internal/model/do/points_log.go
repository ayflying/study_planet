// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsLog is the golang structure of table points_log for DAO operations like Where/Data.
type PointsLog struct {
	g.Meta    `orm:"table:points_log, do:true"`
	Id        any         //
	ChildId   any         //
	Delta     any         //
	Reason    any         //
	CreatedAt *gtime.Time //
}
