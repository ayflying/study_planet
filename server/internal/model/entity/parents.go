// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Parents is the golang structure for table parents.
type Parents struct {
	Id          int64       `json:"id"          orm:"id"            description:""` //
	CasdoorSub  string      `json:"casdoorSub"  orm:"casdoor_sub"   description:""` //
	DisplayName string      `json:"displayName" orm:"display_name"  description:""` //
	Avatar      string      `json:"avatar"      orm:"avatar"        description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:""` //
	LastLoginAt *gtime.Time `json:"lastLoginAt" orm:"last_login_at" description:""` //
}
