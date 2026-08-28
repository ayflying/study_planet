// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Words is the golang structure for table words.
type Words struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	Level     int         `json:"level"     orm:"level"      description:""` //
	Word      string      `json:"word"      orm:"word"       description:""` //
	Meaning   string      `json:"meaning"   orm:"meaning"    description:""` //
	Phonetic  string      `json:"phonetic"  orm:"phonetic"   description:""` //
	Example   string      `json:"example"   orm:"example"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
