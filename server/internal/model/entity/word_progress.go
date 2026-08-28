// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// WordProgress is the golang structure for table word_progress.
type WordProgress struct {
	WordId       int64       `json:"wordId"       orm:"word_id"       description:""` //
	ChildId      int64       `json:"childId"      orm:"child_id"      description:""` //
	Known        int         `json:"known"        orm:"known"         description:""` //
	LastReviewed *gtime.Time `json:"lastReviewed" orm:"last_reviewed" description:""` //
}
