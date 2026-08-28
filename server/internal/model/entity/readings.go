// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Readings is the golang structure for table readings.
type Readings struct {
	Id      int64  `json:"id"      orm:"id"      ` //
	Title   string `json:"title"   orm:"title"   ` //
	Content string `json:"content" orm:"content" ` //
	Level   int    `json:"level"   orm:"level"   ` //
}
