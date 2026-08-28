package v1

import "github.com/gogf/gf/v2/frame/g"

// ImportQuestion 导入题目结构（家长通用导入通道）。
type ImportQuestion struct {
	Subject     string   `json:"subject"`
	Grade       int      `json:"grade"`
	Topic       string   `json:"topic"`
	QType       string   `json:"qtype"`
	Passage     string   `json:"passage,omitempty"`
	Question    string   `json:"question"`
	Options     []string `json:"options"`
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation,omitempty"`
	Difficulty  int      `json:"difficulty"`
	Source      string   `json:"source,omitempty"`
}

// ImportContentReq 通用题目导入（家长鉴权）：按 content_hash 去重，重复自动跳过。
type ImportContentReq struct {
	g.Meta    `path:"/parent/content/import" method:"post" tags:"Content" summary:"导入题目"`
	Questions []ImportQuestion `json:"questions" v:"required"`
}
type ImportContentRes struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
	Total    int `json:"total"`
}

// SubjectStatsReq 内容库统计（家长鉴权）。
type SubjectStatsReq struct {
	g.Meta `path:"/parent/content/stats" method:"get" tags:"Content" summary:"内容库统计"`
}
type SubjectStatsRes struct {
	Total    int       `json:"total"`
	Subjects []Subject `json:"subjects"`
}
