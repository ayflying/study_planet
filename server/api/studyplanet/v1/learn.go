package v1

import "github.com/gogf/gf/v2/frame/g"

// Word 单词卡片。
type Word struct {
	ID        int    `json:"id"`
	Level     int    `json:"level"`
	Word      string `json:"word"`
	Meaning   string `json:"meaning"`
	Phonetic  string `json:"phonetic"`
	Example   string `json:"example"`
	CreatedAt string `json:"created_at"`
}

// ListWordsReq 单词列表（可选按 level 过滤）。
type ListWordsReq struct {
	g.Meta `path:"/words" method:"get" tags:"Learn" summary:"单词列表"`
	Level  string `json:"level" in:"query"`
}
type ListWordsRes []Word

// WordDetailReq 单词详情 + 当前学生掌握状态。
type WordDetailReq struct {
	g.Meta    `path:"/words/:id" method:"get" tags:"Learn" summary:"单词详情"`
	ID        int `in:"path" json:"-"`
	StudentID int `json:"student_id" in:"query"`
}
type WordDetailRes struct {
	Word  Word `json:"word"`
	Known int  `json:"known"`
}

// WordProgressReq 标记单词掌握状态；session_id 传入则走连击+场次计分。
type WordProgressReq struct {
	g.Meta    `path:"/words/:id/progress" method:"post" tags:"Learn" summary:"标记单词掌握"`
	ID        int  `in:"path" json:"-"`
	StudentID int  `json:"student_id" in:"query"`
	Known     bool `json:"known"`
	SessionID int  `json:"session_id"`
}
type WordProgressRes struct {
	OK      bool `json:"ok,omitempty"`
	Known   int  `json:"known"`
	Correct bool `json:"correct,omitempty"`
	// 场次作答（session_id>0）时返回的连击/XP 反馈字段。
	Combo      int `json:"combo,omitempty"`
	BasePoints int `json:"base_points,omitempty"`
	ComboBonus int `json:"combo_bonus,omitempty"`
	Review     int `json:"review,omitempty"`
	XP         int `json:"xp,omitempty"`
}

// Reading 阅读理解短文。
type Reading struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Level   int    `json:"level"`
}

// ReadingQuestion 阅读理解题目。
type ReadingQuestion struct {
	ID        int    `json:"id"`
	ReadingID int    `json:"reading_id"`
	Question  string `json:"question"`
	OptionA   string `json:"option_a"`
	OptionB   string `json:"option_b"`
	OptionC   string `json:"option_c"`
	OptionD   string `json:"option_d"`
	Answer    string `json:"answer"`
}

// ReadingDetailReq 阅读详情 + 题目列表。
type ReadingDetailReq struct {
	g.Meta    `path:"/readings/:id" method:"get" tags:"Learn" summary:"阅读详情"`
	ID        int `in:"path" json:"-"`
	StudentID int `json:"student_id" in:"query"`
}
type ReadingDetailRes struct {
	Reading   Reading           `json:"reading"`
	Questions []ReadingQuestion `json:"questions"`
}

// ReadingAnswerReq 阅读题目作答。
type ReadingAnswerReq struct {
	g.Meta     `path:"/readings/:id/answer" method:"post" tags:"Learn" summary:"阅读作答"`
	ID         int    `in:"path" json:"-"`
	StudentID  int    `json:"student_id" in:"query"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
	SessionID  int    `json:"session_id"`
}
type ReadingAnswerRes struct {
	Correct       bool   `json:"correct"`
	CorrectAnswer string `json:"correct_answer"`
	// 场次作答（session_id>0）时返回的连击/XP 反馈字段。
	Combo      int `json:"combo,omitempty"`
	BasePoints int `json:"base_points,omitempty"`
	ComboBonus int `json:"combo_bonus,omitempty"`
	Review     int `json:"review,omitempty"`
	XP         int `json:"xp,omitempty"`
}

// MathProblem 数学题目。
type MathProblem struct {
	ID          int    `json:"id"`
	Level       int    `json:"level"`
	Type        string `json:"type"`
	Question    string `json:"question"`
	Options     string `json:"options"` // JSON 数组字符串
	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`
}

// ListMathReq 数学题目列表。
type ListMathReq struct {
	g.Meta `path:"/math" method:"get" tags:"Learn" summary:"数学题列表"`
	Level  string `json:"level" in:"query"`
}
type ListMathRes []MathProblem

// MathAnswerReq 数学题作答。
type MathAnswerReq struct {
	g.Meta    `path:"/math/:id/answer" method:"post" tags:"Learn" summary:"数学作答"`
	ID        int    `in:"path" json:"-"`
	StudentID int    `json:"student_id" in:"query"`
	Answer    string `json:"answer"`
	SessionID int    `json:"session_id"`
}
type MathAnswerRes struct {
	Correct     bool   `json:"correct"`
	Explanation string `json:"explanation"`
	Answer      string `json:"answer"`
	// 场次作答（session_id>0）时返回的连击/XP 反馈字段。
	Combo      int `json:"combo,omitempty"`
	BasePoints int `json:"base_points,omitempty"`
	ComboBonus int `json:"combo_bonus,omitempty"`
	Review     int `json:"review,omitempty"`
	XP         int `json:"xp,omitempty"`
}
