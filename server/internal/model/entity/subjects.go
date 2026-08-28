// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Subjects is the golang structure for table subjects.
type Subjects struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	Code      string      `json:"code"      orm:"code"       description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	Icon      string      `json:"icon"      orm:"icon"       description:""` //
	Color     string      `json:"color"     orm:"color"      description:""` //
	MinGrade  int         `json:"minGrade"  orm:"min_grade"  description:""` //
	MaxGrade  int         `json:"maxGrade"  orm:"max_grade"  description:""` //
	Sort      int         `json:"sort"      orm:"sort"       description:""` //
	Enabled   int         `json:"enabled"   orm:"enabled"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
