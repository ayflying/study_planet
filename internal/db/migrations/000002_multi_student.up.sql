-- 000002 多学生模型与家长账号体系
-- 说明：SQLite 方言；切 Postgres/MySQL 时复制到 migrations/postgres/ 调整类型即可。

-- 学生档案：在原 children 表上扩展，老数据自动获得默认值
ALTER TABLE children ADD COLUMN username TEXT NOT NULL DEFAULT '';
ALTER TABLE children ADD COLUMN avatar   TEXT NOT NULL DEFAULT '🚀';
ALTER TABLE children ADD COLUMN grade    INTEGER NOT NULL DEFAULT 5;

-- 用户名可留空（留空则仅用 id 切换）；非空时唯一
CREATE UNIQUE INDEX IF NOT EXISTS idx_children_username ON children(username) WHERE username <> '';

-- 任务与兑换归属到具体学生（老数据归入 1 号学生）
ALTER TABLE tasks       ADD COLUMN child_id INTEGER NOT NULL DEFAULT 1;
ALTER TABLE redemptions ADD COLUMN child_id INTEGER NOT NULL DEFAULT 1;

CREATE INDEX IF NOT EXISTS idx_tasks_child ON tasks(child_id);
CREATE INDEX IF NOT EXISTS idx_redemptions_child ON redemptions(child_id);

-- 家长账号：Casdoor SSO 登录后落库（PIN 登录不依赖本表）
CREATE TABLE IF NOT EXISTS parents (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  casdoor_sub   TEXT NOT NULL UNIQUE,
  display_name  TEXT NOT NULL DEFAULT '',
  avatar        TEXT NOT NULL DEFAULT '',
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_login_at TIMESTAMP
);
