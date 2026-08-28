// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Parents is the golang structure of table parents for DAO operations like Where/Data.
type Parents struct {
	g.Meta      `orm:"table:parents, do:true"`
	Id          interface{} //
	CasdoorSub  interface{} //
	DisplayName interface{} //
	Avatar      interface{} //
	CreatedAt   *gtime.Time //
	LastLoginAt *gtime.Time //
}
