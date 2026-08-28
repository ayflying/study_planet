// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// WordProgress is the golang structure of table word_progress for DAO operations like Where/Data.
type WordProgress struct {
	g.Meta       `orm:"table:word_progress, do:true"`
	WordId       interface{} //
	ChildId      interface{} //
	Known        interface{} //
	LastReviewed *gtime.Time //
}
