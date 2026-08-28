// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tasks is the golang structure for table tasks.
type Tasks struct {
	Id          int64       `json:"id"          orm:"id"           ` //
	Title       string      `json:"title"       orm:"title"        ` //
	Type        string      `json:"type"        orm:"type"         ` //
	DueDate     *gtime.Time `json:"dueDate"     orm:"due_date"     ` //
	Points      int         `json:"points"      orm:"points"       ` //
	Status      string      `json:"status"      orm:"status"       ` //
	ChildId     int64       `json:"childId"     orm:"child_id"     ` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` //
	CompletedAt *gtime.Time `json:"completedAt" orm:"completed_at" ` //
}
