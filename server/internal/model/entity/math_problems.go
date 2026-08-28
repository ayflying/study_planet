// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// MathProblems is the golang structure for table math_problems.
type MathProblems struct {
	Id          int64  `json:"id"          orm:"id"          ` //
	Level       int    `json:"level"       orm:"level"       ` //
	Type        string `json:"type"        orm:"type"        ` //
	Question    string `json:"question"    orm:"question"    ` //
	Options     string `json:"options"     orm:"options"     ` //
	Answer      string `json:"answer"      orm:"answer"      ` //
	Explanation string `json:"explanation" orm:"explanation" ` //
}
