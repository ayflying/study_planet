-- 000004 错题本 + 经验值 + 每周排行榜持久化（MySQL 方言）

CREATE TABLE IF NOT EXISTS wrong_questions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  child_id BIGINT NOT NULL,
  subject VARCHAR(32) NOT NULL,
  ref_id BIGINT NOT NULL,
  wrong_count INT NOT NULL DEFAULT 1,
  resolved TINYINT NOT NULL DEFAULT 0,
  last_wrong_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_reviewed_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_wrong (child_id, subject, ref_id),
  KEY idx_wrong_child (child_id, resolved)
);

ALTER TABLE children ADD COLUMN xp INT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS leaderboard_weekly (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  week_key VARCHAR(16) NOT NULL,
  child_id BIGINT NOT NULL,
  xp INT NOT NULL DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uq_lb (week_key, child_id)
);
