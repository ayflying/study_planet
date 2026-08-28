// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Children is the golang structure of table children for DAO operations like Where/Data.
type Children struct {
	g.Meta    `orm:"table:children, do:true"`
	Id        any         //
	Name      any         //
	Username  any         //
	Avatar    any         //
	Grade     any         //
	CreatedAt *gtime.Time //
	Xp        any         //
}
