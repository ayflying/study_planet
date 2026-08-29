// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Rewards is the golang structure for table rewards.
type Rewards struct {
	Id         int64  `json:"id"         orm:"id"          description:""` //
	Name       string `json:"name"       orm:"name"        description:""` //
	CostPoints int    `json:"costPoints" orm:"cost_points" description:""` //
	Status     string `json:"status"     orm:"status"      description:""` //
	ParentId   *int64 `json:"parentId"   orm:"parent_id"   description:"归属家长，NULL=未归属"` //
}
