-- 000002 回滚：移除多学生与家长账号体系
DROP TABLE IF EXISTS parents;
DROP INDEX IF EXISTS idx_redemptions_child;
DROP INDEX IF EXISTS idx_tasks_child;
DROP INDEX IF EXISTS idx_children_username;
ALTER TABLE redemptions DROP COLUMN child_id;
ALTER TABLE tasks DROP COLUMN child_id;
