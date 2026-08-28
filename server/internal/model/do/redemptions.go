// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Redemptions is the golang structure of table redemptions for DAO operations like Where/Data.
type Redemptions struct {
	g.Meta      `orm:"table:redemptions, do:true"`
	Id          any         //
	RewardId    any         //
	ChildId     any         //
	Status      any         //
	RequestedAt *gtime.Time //
	ConfirmedAt *gtime.Time //
}
