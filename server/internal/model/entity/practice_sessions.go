// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PracticeSessions is the golang structure for table practice_sessions.
type PracticeSessions struct {
	Id         int64       `json:"id"         orm:"id"          ` //
	ChildId    int64       `json:"childId"    orm:"child_id"    ` //
	Subject    string      `json:"subject"    orm:"subject"     ` //
	Level      int         `json:"level"      orm:"level"       ` //
	Total      int         `json:"total"      orm:"total"       ` //
	Correct    int         `json:"correct"    orm:"correct"     ` //
	MaxCombo   int         `json:"maxCombo"   orm:"max_combo"   ` //
	Bonus      int         `json:"bonus"      orm:"bonus"       ` //
	Stars      int         `json:"stars"      orm:"stars"       ` //
	Finished   int         `json:"finished"   orm:"finished"    ` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` //
	FinishedAt *gtime.Time `json:"finishedAt" orm:"finished_at" ` //
}
