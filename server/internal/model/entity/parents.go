// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Parents is the golang structure for table parents.
type Parents struct {
	Id          int64       `json:"id"          orm:"id"            ` //
	CasdoorSub  string      `json:"casdoorSub"  orm:"casdoor_sub"   ` //
	DisplayName string      `json:"displayName" orm:"display_name"  ` //
	Avatar      string      `json:"avatar"      orm:"avatar"        ` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    ` //
	LastLoginAt *gtime.Time `json:"lastLoginAt" orm:"last_login_at" ` //
}
