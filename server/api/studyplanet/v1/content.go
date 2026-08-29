package v1

import "github.com/gogf/gf/v2/frame/g"

// Subject 学科目录条目。
type Subject struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Color    string `json:"color"`
	MinGrade int    `json:"min_grade"`
	MaxGrade int    `json:"max_grade"`
	Sort     int    `json:"sort"`
	Enabled  int    `json:"enabled"`
	Count    int    `json:"count"`
}

// ListSubjectsReq 学科目录；传 grade 只返回该学段开设的学科。
type ListSubjectsReq struct {
	g.Meta `path:"/subjects" method:"get" tags:"Content" summary:"学科目录"`
	Grade  int `json:"grade" in:"query"`
}
type ListSubjectsRes []Subject

// PubQuestion 对外发布的题目（不含 answer，防泄露）。
type PubQuestion struct {
	ID         int      `json:"id"`
	Subject    string   `json:"subject"`
	Grade      int      `json:"grade"`
	Topic      string   `json:"topic"`
	QType      string   `json:"qtype"`
	Passage    string   `json:"passage,omitempty"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	Difficulty int      `json:"difficulty"`
}

// PickQuestionsReq 随机抽题：subject 必填，grade 取该年级相邻难度，limit 默认 5。
type PickQuestionsReq struct {
	g.Meta    `path:"/content/pick" method:"get" tags:"Content" summary:"随机抽题"`
	Subject   string `json:"subject" v:"required" in:"query"`
	Grade     int    `json:"grade" in:"query"`
	Limit     int    `json:"limit" in:"query"`
	StudentID int    `json:"student_id" in:"query"`
}
type PickQuestionsRes []PubQuestion

// ContentItemReq 按 id 取单题（不含答案），错题巩固回取用。
type ContentItemReq struct {
	g.Meta `path:"/content/item" method:"get" tags:"Content" summary:"题目详情"`
	ID     int `json:"id" in:"query" v:"required"`
}
type ContentItemRes PubQuestion

// ContentAnswerReq 统一判分；session_id 传入则复用连击+XP+错题本链路。
type ContentAnswerReq struct {
	g.Meta    `path:"/content/answer" method:"post" tags:"Content" summary:"题目判分"`
	ID        int    `json:"id" v:"required"`
	StudentID int    `json:"student_id" in:"query"`
	Answer    string `json:"answer"`
	SessionID int    `json:"session_id"`
}
type ContentAnswerRes struct {
	Correct     bool   `json:"correct"`
	Answer      string `json:"answer,omitempty"`
	Explanation string `json:"explanation,omitempty"`
	// 场次作答（session_id>0）时返回的连击/XP 反馈字段。
	Combo      int `json:"combo,omitempty"`
	BasePoints int `json:"base_points,omitempty"`
	ComboBonus int `json:"combo_bonus,omitempty"`
	Review     int `json:"review,omitempty"`
	XP         int `json:"xp,omitempty"`
}
