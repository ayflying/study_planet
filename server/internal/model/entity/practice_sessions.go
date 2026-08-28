// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PracticeSessions is the golang structure for table practice_sessions.
type PracticeSessions struct {
	Id         int64       `json:"id"         orm:"id"          description:""` //
	ChildId    int64       `json:"childId"    orm:"child_id"    description:""` //
	Subject    string      `json:"subject"    orm:"subject"     description:""` //
	Level      int         `json:"level"      orm:"level"       description:""` //
	Total      int         `json:"total"      orm:"total"       description:""` //
	Correct    int         `json:"correct"    orm:"correct"     description:""` //
	MaxCombo   int         `json:"maxCombo"   orm:"max_combo"   description:""` //
	Bonus      int         `json:"bonus"      orm:"bonus"       description:""` //
	Stars      int         `json:"stars"      orm:"stars"       description:""` //
	Finished   int         `json:"finished"   orm:"finished"    description:""` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""` //
	FinishedAt *gtime.Time `json:"finishedAt" orm:"finished_at" description:""` //
}
