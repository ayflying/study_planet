// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tasks is the golang structure for table tasks.
type Tasks struct {
	Id          int64       `json:"id"          orm:"id"           description:""` //
	Title       string      `json:"title"       orm:"title"        description:""` //
	Type        string      `json:"type"        orm:"type"         description:""` //
	DueDate     *gtime.Time `json:"dueDate"     orm:"due_date"     description:""` //
	Points      int         `json:"points"      orm:"points"       description:""` //
	Status      string      `json:"status"      orm:"status"       description:""` //
	ChildId     int64       `json:"childId"     orm:"child_id"     description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""` //
	CompletedAt *gtime.Time `json:"completedAt" orm:"completed_at" description:""` //
}
