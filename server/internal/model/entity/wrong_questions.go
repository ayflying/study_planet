// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// WrongQuestions is the golang structure for table wrong_questions.
type WrongQuestions struct {
	Id             int64       `json:"id"             orm:"id"               description:""` //
	ChildId        int64       `json:"childId"        orm:"child_id"         description:""` //
	Subject        string      `json:"subject"        orm:"subject"          description:""` //
	RefId          int64       `json:"refId"          orm:"ref_id"           description:""` //
	WrongCount     int         `json:"wrongCount"     orm:"wrong_count"      description:""` //
	Resolved       int         `json:"resolved"       orm:"resolved"         description:""` //
	LastWrongAt    *gtime.Time `json:"lastWrongAt"    orm:"last_wrong_at"    description:""` //
	LastReviewedAt *gtime.Time `json:"lastReviewedAt" orm:"last_reviewed_at" description:""` //
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""` //
}
