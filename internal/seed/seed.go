package seed

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Run 在空库时写入五年级示例数据（单词/阅读/数学/任务含逾期/奖励/默认孩子/PIN）。
// 幂等：children 已有数据则直接返回。
func Run(db *sqlx.DB, pin string) error {
	var cnt int
	if err := db.Get(&cnt, "SELECT COUNT(*) FROM children"); err != nil {
		return fmt.Errorf("check children: %w", err)
	}
	if cnt > 0 {
		return nil
	}

	if _, err := db.Exec("INSERT INTO children(name) VALUES(?)", "小朋友"); err != nil {
		return err
	}

	// ---- 单词卡片（五年级常用词）----
	words := [][4]string{
		{"because", "因为", "/bɪˈkɒz/", "I stayed home because it rained."},
		{"beautiful", "美丽的", "/ˈbjuːtɪfl/", "She has a beautiful smile."},
		{"favorite", "最喜欢的", "/ˈfeɪvərɪt/", "This is my favorite book."},
		{"usually", "通常", "/ˈjuːʒuəli/", "We usually eat dinner at seven."},
		{"expensive", "昂贵的", "/ɪkˈspensɪv/", "The phone is too expensive."},
		{"dangerous", "危险的", "/ˈdeɪndʒərəs/", "Be careful, it is dangerous."},
		{"library", "图书馆", "/ˈlaɪbrəri/", "I read books in the library."},
		{"Wednesday", "星期三", "/ˈwenzdeɪ/", "We have art on Wednesday."},
		{"dictionary", "字典", "/ˈdɪkʃəneri/", "Use a dictionary to learn words."},
		{"exercise", "练习", "/ˈeksəsaɪz/", "Do exercise every morning."},
	}
	for _, w := range words {
		if _, err := db.Exec(
			"INSERT INTO words(level,word,meaning,phonetic,example) VALUES(?,?,?,?,?)",
			5, w[0], w[1], w[2], w[3],
		); err != nil {
			return err
		}
	}

	// ---- 语文阅读（寓言）----
	tortoise := "乌龟和兔子赛跑，兔子跑得很快，中途骄傲地睡了一觉；乌龟一步一步坚持爬，最终先到终点。故事告诉我们：骄傲使人落后，坚持就是胜利。"
	if _, err := db.Exec("INSERT INTO readings(title,content,level) VALUES(?,?,?)", "龟兔赛跑", tortoise, 5); err != nil {
		return err
	}
	var rid1 int
	if err := db.Get(&rid1, "SELECT id FROM readings WHERE title=?", "龟兔赛跑"); err != nil {
		return err
	}
	for _, q := range [][6]string{
		{"兔子为什么输了比赛？", "它跑得太慢", "它骄傲睡觉", "它迷路了", "它受伤了", "它骄傲睡觉"},
		{"这个故事告诉我们什么？", "要聪明", "坚持就能胜利", "要跑得快", "要睡觉", "坚持就能胜利"},
	} {
		if _, err := db.Exec(
			"INSERT INTO reading_questions(reading_id,question,option_a,option_b,option_c,option_d,answer) VALUES(?,?,?,?,?,?,?)",
			rid1, q[0], q[1], q[2], q[3], q[4], q[5],
		); err != nil {
			return err
		}
	}

	tree := "一个农夫偶然捡到一只撞死在树桩上的兔子，便放下农活天天守着树桩等兔子，结果庄稼荒废，再也没等到兔子。比喻不知变通、妄想不劳而获。"
	if _, err := db.Exec("INSERT INTO readings(title,content,level) VALUES(?,?,?)", "守株待兔", tree, 5); err != nil {
		return err
	}
	var rid2 int
	if err := db.Get(&rid2, "SELECT id FROM readings WHERE title=?", "守株待兔"); err != nil {
		return err
	}
	for _, q := range [][6]string{
		{"农夫为什么再也等不到兔子？", "兔子变聪明了", "兔子不会总撞树", "他搬家了", "天气变了", "兔子不会总撞树"},
		{"这个成语讽刺了哪种人？", "勤劳的人", "妄想不劳而获的人", "勇敢的人", "聪明的人", "妄想不劳而获的人"},
	} {
		if _, err := db.Exec(
			"INSERT INTO reading_questions(reading_id,question,option_a,option_b,option_c,option_d,answer) VALUES(?,?,?,?,?,?,?)",
			rid2, q[0], q[1], q[2], q[3], q[4], q[5],
		); err != nil {
			return err
		}
	}

	// ---- 数学题目（小数/面积/方程/分数）----
	for _, m := range [][7]string{
		{"计算：3.25 + 1.75 = ?", "4.0", "5.0", "5.5", "4.5", "5.0", "小数相加：3.25+1.75=5.00"},
		{"一个长方形长6cm、宽4cm，面积是？", "10", "20", "24", "12", "24", "长方形面积=长×宽=6×4=24"},
		{"解方程：2x + 4 = 14，x = ?", "4", "5", "6", "7", "5", "2x=10，x=5"},
		{"计算：1/2 + 1/4 = ?", "3/4", "2/3", "1/2", "1", "3/4", "通分：2/4+1/4=3/4"},
	} {
		opts := `["` + m[1] + `","` + m[2] + `","` + m[3] + `","` + m[4] + `"]`
		if _, err := db.Exec(
			"INSERT INTO math_problems(level,type,question,options,answer,explanation) VALUES(?,?,?,?,?,?)",
			5, "计算", m[0], opts, m[5], m[6],
		); err != nil {
			return err
		}
	}

	// ---- 每日任务（含 1 条逾期）----
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	for _, t := range [][4]string{
		{"背诵10个英语单词", "单词", today, "10"},
		{"完成数学练习一页", "数学", today, "15"},
		{"读一篇语文课文", "语文", yesterday, "10"},
	} {
		pts, _ := strconv.Atoi(t[3])
		if _, err := db.Exec(
			"INSERT INTO tasks(title,type,due_date,points,status) VALUES(?,?,?,?,'pending')",
			t[0], t[1], t[2], pts,
		); err != nil {
			return err
		}
	}

	// ---- 积分奖励 ----
	for _, rw := range [][2]string{
		{"看动画片30分钟", "50"},
		{"买一本喜欢的绘本", "80"},
		{"周末去公园玩", "100"},
	} {
		cost, _ := strconv.Atoi(rw[1])
		if _, err := db.Exec("INSERT INTO rewards(name,cost_points,status) VALUES(?,?,'active')", rw[0], cost); err != nil {
			return err
		}
	}

	// ---- 家长 PIN（bcrypt 哈希存入 settings）----
	hash, err := bcrypt.GenerateFromPassword([]byte(pin), 10)
	if err != nil {
		return err
	}
	if _, err := db.Exec(
		"INSERT INTO settings(key,value) VALUES('parent_pin',?) ON CONFLICT(key) DO UPDATE SET value=?",
		string(hash), string(hash),
	); err != nil {
		return err
	}

	return nil
}
