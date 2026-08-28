# 数据库切换 MySQL 与自动迁移

## 连接配置（.env，不入库）

```text
DB_DRIVER=mysql
DB_DSN=study_planet:***@tcp(100.66.1.1:3306)/study_planet?charset=utf8mb4&parseTime=true&loc=Local
```

## 迁移机制

- `server/internal/db/db.go` 自研轻量迁移器：启动时自动建 `schema_migrations` 版本表，空库全量建表，已有库按版本增量执行；检测到版本回退（库版本大于脚本版本）直接报错拒绝启动，避免静默破坏数据。
- 迁移脚本按方言分目录内嵌：`server/internal/db/migrations/sqlite/`（原 000001-000003）与 `server/internal/db/migrations/mysql/`（000001 全量建表，含 children/words/word_progress/readings/reading_questions/math_problems/tasks/points_log/rewards/redemptions/settings/parents/practice_sessions/session_answers）。
- 原计划使用 golang-migrate，但 v4.19.1 的 mysql 驱动不支持 multiStatement 且会留下脏版本记录（Error 1064 后标记 dirty=true 导致后续迁移被跳过），故改为逐语句执行的自研迁移器，并需手动清理遗留 `schema_migrations`。

## MySQL 兼容性修复

- `settings` 表主键列名 `key` 是 MySQL 保留字：所有 SQL 加反引号（`` `key` ``），包括 seed 写 PIN、SetPin、ParentLogin 查询。
- `children.username` 改为 nullable + 唯一索引：空用户名写 NULL（MySQL 唯一索引允许多个 NULL），非空用户名唯一；Go 侧用 `sql.NullString` 写入，查询用 `COALESCE(username,'') AS username` 输出字符串。SQLite 分支不受影响（partial index 语义保留）。
- upsert 语句 `ON CONFLICT ... DO UPDATE` 改为 MySQL 的 `ON DUPLICATE KEY UPDATE`（word_progress、parents、settings 三处）。
- `CreateStudent` 用 `res.LastInsertId()` 替代 SQLite 专有的 `last_insert_rowid()`。
- 启动时仅在 SQLite 驱动下创建 DSN 目录（MySQL DSN 含 `tcp(...)` 会被误当成路径）。

## 真实验证结果（本地连 100.66.1.1 MySQL 9.7.0）

- 空库启动 → 迁移全量建 15 张表 + 种子数据（默认学生、10 单词、2 阅读、4 数学、3 任务、3 奖励、PIN）。
- PIN 登录 1234 → 拿到 token；创建学生（空用户名与 kid2 均成功）→ 学生列表输出正常 username 字符串。
- 创建练习场次、单词 progress 上报、积分累计（today_earned=5, total=5）全部通过。
- `go vet` / `go test` 通过；Dockerfile 默认 SQLite，MySQL 由 .env 覆盖。
