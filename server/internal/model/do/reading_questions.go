// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ReadingQuestions is the golang structure of table reading_questions for DAO operations like Where/Data.
type ReadingQuestions struct {
	g.Meta    `orm:"table:reading_questions, do:true"`
	Id        any //
	ReadingId any //
	Question  any //
	OptionA   any //
	OptionB   any //
	OptionC   any //
	OptionD   any //
	Answer    any //
}
