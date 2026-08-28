// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Redemptions is the golang structure for table redemptions.
type Redemptions struct {
	Id          int64       `json:"id"          orm:"id"           description:""` //
	RewardId    int64       `json:"rewardId"    orm:"reward_id"    description:""` //
	ChildId     int64       `json:"childId"     orm:"child_id"     description:""` //
	Status      string      `json:"status"      orm:"status"       description:""` //
	RequestedAt *gtime.Time `json:"requestedAt" orm:"requested_at" description:""` //
	ConfirmedAt *gtime.Time `json:"confirmedAt" orm:"confirmed_at" description:""` //
}
