// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Children is the golang structure for table children.
type Children struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Username  string      `json:"username"  orm:"username"   description:""` //
	Avatar    string      `json:"avatar"    orm:"avatar"     description:""` //
	Grade     int         `json:"grade"     orm:"grade"      description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	Xp        int         `json:"xp"        orm:"xp"         description:""` //
	ParentId  *int64      `json:"parentId"  orm:"parent_id"  description:"归属家长，NULL=未归属"` //
}
