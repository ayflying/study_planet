// Package contentlib 动态内容库：学科目录 + 统一题库。
// 学习内容全部存于数据库，新增资料只需调用导入接口/导入工具，无需改源码。
package contentlib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"context"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// Subject 学科目录条目。
type Subject struct {
	Code     string `db:"code" json:"code"`
	Name     string `db:"name" json:"name"`
	Icon     string `db:"icon" json:"icon"`
	Color    string `db:"color" json:"color"`
	MinGrade int    `db:"min_grade" json:"min_grade"`
	MaxGrade int    `db:"max_grade" json:"max_grade"`
	Sort     int    `db:"sort" json:"sort"`
	Enabled  int    `db:"enabled" json:"enabled"`
}

// Question 统一题目结构（含阅读短文 passage，可空）。
type Question struct {
	Subject     string   `db:"subject" json:"subject"`
	Grade       int      `db:"grade" json:"grade"`
	Topic       string   `db:"topic" json:"topic"`
	QType       string   `db:"qtype" json:"qtype"`
	Passage     string   `db:"passage" json:"passage,omitempty"`
	Question    string   `db:"question" json:"question"`
	Options     []string `db:"-" json:"options"`
	OptionsJSON string   `db:"options" json:"-"`
	Answer      string   `db:"answer" json:"answer"`
	Explanation string   `db:"explanation" json:"explanation,omitempty"`
	Difficulty  int      `db:"difficulty" json:"difficulty"`
	Source      string   `db:"source" json:"source,omitempty"`
}

// builtinSubjects 内置学科目录（小学 1-6 / 初中 7-9）。
var builtinSubjects = []Subject{
	{"english", "英语", "Aa", "#4a90d9", 1, 9, 1, 1},
	{"chinese", "语文", "文", "#e67e22", 1, 9, 2, 1},
	{"math", "数学", "∑", "#27ae60", 1, 9, 3, 1},
	{"physics", "物理", "⚛", "#8e44ad", 8, 9, 4, 1},
	{"chemistry", "化学", "⚗", "#16a085", 9, 9, 5, 1},
	{"biology", "生物", "🧬", "#c0392b", 7, 9, 6, 1},
	{"history", "历史", "📜", "#d4a017", 7, 9, 7, 1},
	{"geography", "地理", "🗺", "#2c7bb6", 7, 9, 8, 1},
}

// HashQuestion 计算题目内容指纹（subject+question+answer 去重）。
func HashQuestion(q *Question) string {
	h := md5.Sum([]byte(q.Subject + "\x1f" + strings.TrimSpace(q.Question) + "\x1f" + strings.TrimSpace(q.Answer)))
	return hex.EncodeToString(h[:])
}

// UpsertSubjects 写入内置学科目录（幂等，按 code 更新展示信息）。
func UpsertSubjects(ctx context.Context, db gdb.DB) error {
	for _, s := range builtinSubjects {
		if _, err := db.Exec(ctx,
			`INSERT INTO subjects(code,name,icon,color,min_grade,max_grade,sort,enabled) VALUES(?,?,?,?,?,?,?,?)
			 ON DUPLICATE KEY UPDATE name=VALUES(name), icon=VALUES(icon), color=VALUES(color), min_grade=VALUES(min_grade), max_grade=VALUES(max_grade), sort=VALUES(sort)`,
			s.Code, s.Name, s.Icon, s.Color, s.MinGrade, s.MaxGrade, s.Sort, s.Enabled,
		); err != nil {
			return fmt.Errorf("写入学科 %s 失败: %w", s.Code, err)
		}
	}
	return nil
}

// ImportQuestions 批量导入题目：按 content_hash 去重（已存在跳过），返回导入条数。
// 这是唯一的题目入库通道：采集脚本/生成器都输出 []Question 后调用本函数。
func ImportQuestions(ctx context.Context, db gdb.DB, qs []Question) (imported int, skipped int, err error) {
	for i := range qs {
		q := &qs[i]
		if err := NormalizeQuestion(q); err != nil {
			return imported, skipped, fmt.Errorf("第 %d 题: %w", i+1, err)
		}
		hash := HashQuestion(q)
		res, err := db.Exec(ctx,
			`INSERT IGNORE INTO questions(subject,grade,topic,qtype,passage,question,options,answer,explanation,difficulty,source,content_hash,enabled)
			 VALUES(?,?,?,?,?,?,?,?,?,?,?,?,1)`,
			q.Subject, q.Grade, q.Topic, q.QType, nullIfEmpty(q.Passage), q.Question, q.OptionsJSON, q.Answer, nullIfEmpty(q.Explanation), q.Difficulty, q.Source, hash,
		)
		if err != nil {
			return imported, skipped, fmt.Errorf("导入第 %d 题失败: %w", i+1, err)
		}
		n, _ := res.RowsAffected()
		if n > 0 {
			imported++
		} else {
			skipped++
		}
	}
	return imported, skipped, nil
}

// NormalizeQuestion 校验并补全题目字段，把 Options 序列化为 JSON（导入通道对外暴露）。
func NormalizeQuestion(q *Question) error {
	q.Subject = strings.TrimSpace(q.Subject)
	if q.Subject == "" {
		return fmt.Errorf("subject 不能为空")
	}
	if strings.TrimSpace(q.Question) == "" {
		return fmt.Errorf("question 不能为空")
	}
	if strings.TrimSpace(q.Answer) == "" {
		return fmt.Errorf("answer 不能为空")
	}
	if q.Grade < 1 || q.Grade > 9 {
		q.Grade = 1
	}
	if q.QType == "" {
		q.QType = "choice"
	}
	if q.Difficulty <= 0 {
		q.Difficulty = 1
	}
	if len(q.Options) > 0 {
		// 统一为字符串数组
		strs := make([]string, len(q.Options))
		for i, o := range q.Options {
			strs[i] = fmt.Sprintf("%v", o)
		}
		q.Options = strs
	}
	if len(q.Options) == 0 && q.QType == "choice" {
		return fmt.Errorf("选择题 options 不能为空: %s", q.Question)
	}
	// 序列化选项为 JSON（入库统一格式）
	if b, err := json.Marshal(q.Options); err == nil {
		q.OptionsJSON = string(b)
	} else {
		return fmt.Errorf("options 序列化失败: %w", err)
	}
	return nil
}

func nullIfEmpty(s string) interface{} {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

// ListSubjects 学科列表（enabled=1 按顺序）。
func ListSubjects(ctx context.Context, db gdb.DB) ([]Subject, error) {
	var ss []Subject
	if err := db.Model("subjects").Ctx(ctx).Where("enabled", 1).Order("sort", "id").Scan(&ss); err != nil {
		return nil, err
	}
	if ss == nil {
		ss = []Subject{}
	}
	return ss, nil
}

// CountBySubject 每个学科的题目数（按学段），用于前端展示。
func CountBySubject(ctx context.Context, db gdb.DB) (map[string]int, error) {
	all, err := db.GetAll(ctx, "SELECT subject, COUNT(*) AS cnt FROM questions WHERE enabled=1 GROUP BY subject")
	if err != nil {
		return nil, err
	}
	m := make(map[string]int, len(all))
	for _, r := range all {
		m[r["subject"].String()] = r["cnt"].Int()
	}
	return m, nil
}

// SubjectExists 判断学科 code 是否已启用（含题目为 0 的科目）。
func SubjectExists(ctx context.Context, db gdb.DB, code string) (bool, error) {
	n, err := db.Model("subjects").Ctx(ctx).Where("code", code).Where("enabled", 1).Count()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// SortStrings 工具：给选项排序（导出用）。
func SortStrings(s []string) { sort.Strings(s) }
