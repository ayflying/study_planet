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
	Id        interface{} //
	ReadingId interface{} //
	Question  interface{} //
	OptionA   interface{} //
	OptionB   interface{} //
	OptionC   interface{} //
	OptionD   interface{} //
	Answer    interface{} //
}
