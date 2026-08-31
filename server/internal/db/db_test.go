package db

import "testing"

// 验证版本号冲突检测：两个同版本号文件必须报错（历史踩坑：000008 双文件导致其中一个永不执行）。
func TestLoadMigrationsDuplicateVersion(t *testing.T) {
	// 使用嵌入的真实迁移目录：当前应为合法（无重复版本号）
	ups, err := loadMigrations("mysql")
	if err != nil {
		t.Fatalf("loadMigrations(mysql) 失败: %v", err)
	}
	if len(ups) == 0 {
		t.Fatal("未加载到任何迁移")
	}
	seen := map[int]string{}
	for _, m := range ups {
		if prev, dup := seen[m.version]; dup {
			t.Fatalf("迁移版本 %06d 重复: %s 与 %s", m.version, prev, m.name)
		}
		seen[m.version] = m.name
	}
}

// 验证 000009 迁移脚本不再包含 MySQL 8 不支持的 ADD COLUMN IF NOT EXISTS 语法。
func TestNoMariaDBOnlySyntax(t *testing.T) {
	ups, err := loadMigrations("mysql")
	if err != nil {
		t.Fatalf("loadMigrations(mysql) 失败: %v", err)
	}
	for _, m := range ups {
		stmts := splitStatements(m.upSQL)
		for _, s := range stmts {
			if alterAddColRe.MatchString(s) && containsIfNotExists(s) {
				t.Errorf("迁移 %s 含 MySQL 8 不支持的 IF NOT EXISTS 语法:\n%s", m.name, s)
			}
		}
	}
}

func containsIfNotExists(s string) bool {
	target := "IF NOT EXISTS"
	for i := 0; i+len(target) <= len(s); i++ {
		if s[i:i+len(target)] == target {
			return true
		}
	}
	return false
}

// 验证 statementCovered 的正则匹配：仅 ALTER TABLE ... ADD COLUMN 形态可判定跳过。
func TestAlterAddColRe(t *testing.T) {
	cases := []struct {
		stmt  string
		match bool
		table string
		col   string
	}{
		{"ALTER TABLE pets ADD COLUMN food_apple INT NOT NULL DEFAULT 0", true, "pets", "food_apple"},
		{"ALTER TABLE `pets` ADD COLUMN `stars_spent` INT NOT NULL DEFAULT 0", true, "pets", "stars_spent"},
		{"ALTER TABLE pets ADD COLUMN IF NOT EXISTS food_apple INT NOT NULL DEFAULT 0", true, "pets", "food_apple"},
		{"alter table pets\n add column stars_spent int not null default 0", true, "pets", "stars_spent"},
		{"ALTER TABLE pets DROP COLUMN food_apple", false, "", ""},
		{"CREATE TABLE t (id INT)", false, "", ""},
		{"ALTER TABLE pets ADD KEY idx_x (col)", false, "", ""},
	}
	for _, c := range cases {
		m := alterAddColRe.FindStringSubmatch(c.stmt)
		if (m != nil) != c.match {
			t.Errorf("stmt=%q 期望 match=%v 实际 %v", c.stmt, c.match, m != nil)
			continue
		}
		if m != nil && (m[1] != c.table || m[2] != c.col) {
			t.Errorf("stmt=%q 期望 table=%s col=%s 实际 %s/%s", c.stmt, c.table, c.col, m[1], m[2])
		}
	}
}

// 验证 splitStatements 保留完整语句且剔除注释行。
func TestSplitStatements(t *testing.T) {
	sqlText := "-- 注释行\nALTER TABLE a ADD COLUMN x INT;\r\n-- 另一条注释\nALTER TABLE b ADD COLUMN y INT;\n"
	stmts := splitStatements(sqlText)
	if len(stmts) != 2 {
		t.Fatalf("期望 2 条语句，实际 %d: %q", len(stmts), stmts)
	}
	if stmts[0] != "ALTER TABLE a ADD COLUMN x INT" || stmts[1] != "ALTER TABLE b ADD COLUMN y INT" {
		t.Errorf("语句切分结果异常: %q", stmts)
	}
}
