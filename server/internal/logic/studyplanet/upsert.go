package studyplanet

import "github.com/jmoiron/sqlx"

// upsertExec 方言自适应 upsert：同一张表按连接驱动选择 MySQL 或 SQLite 语法。
//   - columns/vals: 插入列与值
//   - conflictCols: 冲突判定列（MySQL 用主键/唯一索引推断；SQLite 显式声明）
//   - updates:      冲突时要更新的 "col=excluded.col" 列表
//
// 语义统一为「存在则覆盖更新列」。需要累加语义（xp= xp+...）的场景不适用本函数，请单独写。
func upsertExec(db *sqlx.DB, table string, columns []string, vals []interface{}, conflictCols []string, updates []string) error {
	var q string
	if db.DriverName() == "mysql" {
		set := ""
		for i, u := range updates {
			if i > 0 {
				set += ", "
			}
			set += u + "=VALUES(" + u + ")"
		}
		q = "INSERT INTO " + table + "(" + joinCols(columns) + ") VALUES(" + placeholders(len(vals)) + ") ON DUPLICATE KEY UPDATE " + set
	} else {
		set := ""
		for i, u := range updates {
			if i > 0 {
				set += ", "
			}
			set += u + "=excluded." + u
		}
		q = "INSERT INTO " + table + "(" + joinCols(columns) + ") VALUES(" + placeholders(len(vals)) + ") ON CONFLICT(" + joinCols(conflictCols) + ") DO UPDATE SET " + set
	}
	_, err := db.Exec(q, vals...)
	return err
}

func joinCols(cols []string) string {
	s := ""
	for i, c := range cols {
		if i > 0 {
			s += ", "
		}
		s += c
	}
	return s
}

func placeholders(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			s += ","
		}
		s += "?"
	}
	return s
}
