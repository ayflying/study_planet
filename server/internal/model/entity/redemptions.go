// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Redemptions is the golang structure for table redemptions.
type Redemptions struct {
	Id          int64       `json:"id"          orm:"id"           ` //
	RewardId    int64       `json:"rewardId"    orm:"reward_id"    ` //
	ChildId     int64       `json:"childId"     orm:"child_id"     ` //
	Status      string      `json:"status"      orm:"status"       ` //
	RequestedAt *gtime.Time `json:"requestedAt" orm:"requested_at" ` //
	ConfirmedAt *gtime.Time `json:"confirmedAt" orm:"confirmed_at" ` //
}
