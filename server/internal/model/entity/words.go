// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Words is the golang structure for table words.
type Words struct {
	Id        int64       `json:"id"        orm:"id"         ` //
	Level     int         `json:"level"     orm:"level"      ` //
	Word      string      `json:"word"      orm:"word"       ` //
	Meaning   string      `json:"meaning"   orm:"meaning"    ` //
	Phonetic  string      `json:"phonetic"  orm:"phonetic"   ` //
	Example   string      `json:"example"   orm:"example"    ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
}
