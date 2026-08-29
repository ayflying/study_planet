-- 000008 家长数据隔离：孩子与奖励归属到具体家长账号（parents.id，SQLite 方言）
-- 约定：parent_id 为 NULL 表示未归属（由第一个登录的 Casdoor 家长自动接管）。

ALTER TABLE children ADD COLUMN parent_id INTEGER NULL DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_children_parent ON children(parent_id);

ALTER TABLE rewards ADD COLUMN parent_id INTEGER NULL DEFAULT NULL;
CREATE INDEX IF NOT EXISTS idx_rewards_parent ON rewards(parent_id);
