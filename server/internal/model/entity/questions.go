// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure for table questions.
type Questions struct {
	Id          int64       `json:"id"          orm:"id"           description:""` //
	Subject     string      `json:"subject"     orm:"subject"      description:""` //
	Grade       int         `json:"grade"       orm:"grade"        description:""` //
	Topic       string      `json:"topic"       orm:"topic"        description:""` //
	Qtype       string      `json:"qtype"       orm:"qtype"        description:""` //
	Passage     string      `json:"passage"     orm:"passage"      description:""` //
	Question    string      `json:"question"    orm:"question"     description:""` //
	Options     string      `json:"options"     orm:"options"      description:""` //
	Answer      string      `json:"answer"      orm:"answer"       description:""` //
	Explanation string      `json:"explanation" orm:"explanation"  description:""` //
	Difficulty  int         `json:"difficulty"  orm:"difficulty"   description:""` //
	Source      string      `json:"source"      orm:"source"       description:""` //
	ContentHash string      `json:"contentHash" orm:"content_hash" description:""` //
	Enabled     int         `json:"enabled"     orm:"enabled"      description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""` //
}
