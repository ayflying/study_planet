// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// MathProblems is the golang structure for table math_problems.
type MathProblems struct {
	Id          int64  `json:"id"          orm:"id"          description:""` //
	Level       int    `json:"level"       orm:"level"       description:""` //
	Type        string `json:"type"        orm:"type"        description:""` //
	Question    string `json:"question"    orm:"question"    description:""` //
	Options     string `json:"options"     orm:"options"     description:""` //
	Answer      string `json:"answer"      orm:"answer"      description:""` //
	Explanation string `json:"explanation" orm:"explanation" description:""` //
}
