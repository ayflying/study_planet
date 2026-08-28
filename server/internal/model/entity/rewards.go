// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Rewards is the golang structure for table rewards.
type Rewards struct {
	Id         int64  `json:"id"         orm:"id"          ` //
	Name       string `json:"name"       orm:"name"        ` //
	CostPoints int    `json:"costPoints" orm:"cost_points" ` //
	Status     string `json:"status"     orm:"status"      ` //
}
