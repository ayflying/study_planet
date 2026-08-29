-- 000006 家长数据隔离：孩子与奖励归属到具体家长账号（parents.id）
-- 约定：parent_id 为 NULL 表示未归属（由第一个登录的 Casdoor 家长自动接管）。

ALTER TABLE children
  ADD COLUMN parent_id BIGINT NULL DEFAULT NULL,
  ADD KEY idx_children_parent (parent_id);

ALTER TABLE rewards
  ADD COLUMN parent_id BIGINT NULL DEFAULT NULL,
  ADD KEY idx_rewards_parent (parent_id);
