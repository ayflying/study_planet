// Package seedcontent 动态内容库启动引导：同步学科目录，空库时导入内置全科题。
// 独立包避免 contentlib ↔ contentgen 循环依赖。
package seedcontent

import (
	"log"

	"github.com/jmoiron/sqlx"

	"studyplanet/internal/contentgen"
	"studyplanet/internal/contentlib"
)

// Run 启动时同步内置学科目录；内置题库只在 questions 表为空时导入一次，
// 之后以数据库为准（管理员可通过导入接口增补，不会被启动覆盖）。
func Run(db *sqlx.DB) error {
	if err := contentlib.UpsertSubjects(db); err != nil {
		return err
	}
	var cnt int
	if err := db.Get(&cnt, "SELECT COUNT(*) FROM questions"); err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}
	qs := contentgen.Generate()
	imported, skipped, err := contentlib.ImportQuestions(db, qs)
	if err != nil {
		return err
	}
	log.Printf("内容库：内置题库首次导入完成，共 %d 题（新增 %d，跳过重复 %d）", len(qs), imported, skipped)
	return nil
}
