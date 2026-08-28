// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PointsLog is the golang structure for table points_log.
type PointsLog struct {
	Id        int64       `json:"id"        orm:"id"         description:""` //
	ChildId   int64       `json:"childId"   orm:"child_id"   description:""` //
	Delta     int         `json:"delta"     orm:"delta"      description:""` //
	Reason    string      `json:"reason"    orm:"reason"     description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
