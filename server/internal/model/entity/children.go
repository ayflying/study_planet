// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Children is the golang structure for table children.
type Children struct {
	Id        int64       `json:"id"        orm:"id"         ` //
	Name      string      `json:"name"      orm:"name"       ` //
	Username  string      `json:"username"  orm:"username"   ` //
	Avatar    string      `json:"avatar"    orm:"avatar"     ` //
	Grade     int         `json:"grade"     orm:"grade"      ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
}
