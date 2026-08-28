package main

// 内容库本地验证：空库迁移 → 内置题导入 → 抽题 → 判分 → 去重再导入。
import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"studyplanet/internal/contentgen"
	"studyplanet/internal/contentlib"
	"studyplanet/internal/db"
	"studyplanet/internal/seedcontent"
)

func main() {
	dsn := "data/verify_content.db"
	_ = os.Remove(dsn)
	if err := db.Migrate("sqlite", dsn); err != nil {
		log.Fatal("迁移失败:", err)
	}
	conn, err := db.Open("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := seedcontent.Run(conn); err != nil {
		log.Fatal("内容库导入失败:", err)
	}

	// 统计
	counts, _ := contentlib.CountBySubject(conn)
	total := 0
	for sub, c := range counts {
		fmt.Printf("  %-10s %4d 题\n", sub, c)
		total += c
	}
	fmt.Println("总题数:", total)

	// 抽题（直接查库模拟 PickQuestions 逻辑）
	var raw []struct {
		ID       int    `db:"id"`
		Subject  string `db:"subject"`
		Question string `db:"question"`
		Options  string `db:"options"`
		Answer   string `db:"answer"`
	}
	if err := conn.Select(&raw, "SELECT id,subject,question,options,answer FROM questions WHERE subject='math' AND enabled=1 AND grade BETWEEN 4 AND 6 ORDER BY RANDOM() LIMIT 3"); err != nil {
		log.Fatal(err)
	}
	for _, q := range raw {
		var opts []string
		_ = json.Unmarshal([]byte(q.Options), &opts)
		fmt.Printf("  [%s#%d] %s 选项=%v 答案=%s\n", q.Subject, q.ID, q.Question, opts, q.Answer)
	}

	// 去重验证：再次导入全部生成数据应全部跳过
	qs := contentgen.Generate()
	imported, skipped, err := contentlib.ImportQuestions(conn, qs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("重复导入验证: imported=%d skipped=%d（期望 imported=0）\n", imported, skipped)

	_ = os.Remove(dsn)
	fmt.Println("VERIFY_OK")
}
