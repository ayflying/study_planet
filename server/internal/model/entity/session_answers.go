// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SessionAnswers is the golang structure for table session_answers.
type SessionAnswers struct {
	Id         int64       `json:"id"         orm:"id"          ` //
	SessionId  int64       `json:"sessionId"  orm:"session_id"  ` //
	RefId      int64       `json:"refId"      orm:"ref_id"      ` //
	Correct    int         `json:"correct"    orm:"correct"     ` //
	Points     int         `json:"points"     orm:"points"      ` //
	Combo      int         `json:"combo"      orm:"combo"       ` //
	AnsweredAt *gtime.Time `json:"answeredAt" orm:"answered_at" ` //
}
