-- 000004 错题本 + 经验值 + 每周排行榜持久化（SQLite 方言）

-- 错题本：答错登记，巩固答对后标记 resolved；同一题重复答错重新激活并累计次数
CREATE TABLE IF NOT EXISTS wrong_questions (
  id               INTEGER PRIMARY KEY AUTOINCREMENT,
  child_id         INTEGER NOT NULL,
  subject          TEXT NOT NULL,             -- words | reading | math
  ref_id           INTEGER NOT NULL,          -- word / reading_question / math_problem 的 id
  wrong_count      INTEGER NOT NULL DEFAULT 1,
  resolved         INTEGER NOT NULL DEFAULT 0, -- 1 = 巩固练习中已答对
  last_wrong_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_reviewed_at TIMESTAMP,
  created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_wrong ON wrong_questions(child_id, subject, ref_id);
CREATE INDEX IF NOT EXISTS idx_wrong_child ON wrong_questions(child_id, resolved);

-- 经验值累计（children.xp，答对/奖励时同步累加）
ALTER TABLE children ADD COLUMN xp INTEGER NOT NULL DEFAULT 0;

-- 每周排行榜持久化表（Redis ZSET 每小时快照落库，防 Redis 数据丢失）
CREATE TABLE IF NOT EXISTS leaderboard_weekly (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  week_key   TEXT NOT NULL,                  -- 例：2026W35（ISO 周，周一为一周开始）
  child_id   INTEGER NOT NULL,
  xp         INTEGER NOT NULL DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_lb ON leaderboard_weekly(week_key, child_id);
