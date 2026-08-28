// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ReadingQuestions is the golang structure for table reading_questions.
type ReadingQuestions struct {
	Id        int64  `json:"id"        orm:"id"         description:""` //
	ReadingId int64  `json:"readingId" orm:"reading_id" description:""` //
	Question  string `json:"question"  orm:"question"   description:""` //
	OptionA   string `json:"optionA"   orm:"option_a"   description:""` //
	OptionB   string `json:"optionB"   orm:"option_b"   description:""` //
	OptionC   string `json:"optionC"   orm:"option_c"   description:""` //
	OptionD   string `json:"optionD"   orm:"option_d"   description:""` //
	Answer    string `json:"answer"    orm:"answer"     description:""` //
}
