// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Settings is the golang structure of table settings for DAO operations like Where/Data.
type Settings struct {
	g.Meta `orm:"table:settings, do:true"`
	Key    any //
	Value  any //
}
