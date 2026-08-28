// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ReadingQuestions is the golang structure for table reading_questions.
type ReadingQuestions struct {
	Id        int64  `json:"id"        orm:"id"         ` //
	ReadingId int64  `json:"readingId" orm:"reading_id" ` //
	Question  string `json:"question"  orm:"question"   ` //
	OptionA   string `json:"optionA"   orm:"option_a"   ` //
	OptionB   string `json:"optionB"   orm:"option_b"   ` //
	OptionC   string `json:"optionC"   orm:"option_c"   ` //
	OptionD   string `json:"optionD"   orm:"option_d"   ` //
	Answer    string `json:"answer"    orm:"answer"     ` //
}
