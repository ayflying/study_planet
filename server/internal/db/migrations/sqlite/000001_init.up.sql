-- 000001 初始化：五年级学习闯关台数据表
-- 说明：当前使用 SQLite 方言；日后升级到 Postgres/MySQL 时，迁移版本机制通用，
-- 只需把本文件复制为 migrations/postgres/ 并调整自增/时间戳类型即可（详见 README）。

CREATE TABLE IF NOT EXISTS children (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  name       TEXT NOT NULL DEFAULT '小朋友',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS words (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  level      INTEGER NOT NULL DEFAULT 1,
  word       TEXT NOT NULL,
  meaning    TEXT NOT NULL DEFAULT '',
  phonetic   TEXT NOT NULL DEFAULT '',
  example    TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS word_progress (
  word_id       INTEGER NOT NULL,
  child_id      INTEGER NOT NULL,
  known         INTEGER NOT NULL DEFAULT 0,
  last_reviewed TIMESTAMP,
  PRIMARY KEY (word_id, child_id)
);

CREATE TABLE IF NOT EXISTS readings (
  id      INTEGER PRIMARY KEY AUTOINCREMENT,
  title   TEXT NOT NULL,
  content TEXT NOT NULL,
  level   INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS reading_questions (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  reading_id INTEGER NOT NULL,
  question   TEXT NOT NULL,
  option_a   TEXT NOT NULL DEFAULT '',
  option_b   TEXT NOT NULL DEFAULT '',
  option_c   TEXT NOT NULL DEFAULT '',
  option_d   TEXT NOT NULL DEFAULT '',
  answer     TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS math_problems (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  level       INTEGER NOT NULL DEFAULT 1,
  type        TEXT NOT NULL DEFAULT '',
  question    TEXT NOT NULL,
  options     TEXT NOT NULL DEFAULT '[]',
  answer      TEXT NOT NULL,
  explanation TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS tasks (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  title       TEXT NOT NULL,
  type        TEXT NOT NULL DEFAULT '',
  due_date    DATE,
  points      INTEGER NOT NULL DEFAULT 0,
  status      TEXT NOT NULL DEFAULT 'pending',
  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS points_log (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  child_id   INTEGER NOT NULL,
  delta      INTEGER NOT NULL DEFAULT 0,
  reason     TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rewards (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL,
  cost_points INTEGER NOT NULL DEFAULT 0,
  status      TEXT NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS redemptions (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  reward_id    INTEGER NOT NULL,
  status       TEXT NOT NULL DEFAULT 'pending',
  requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  confirmed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS settings (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_words_level ON words(level);
CREATE INDEX IF NOT EXISTS idx_reading_q ON reading_questions(reading_id);
CREATE INDEX IF NOT EXISTS idx_math_level ON math_problems(level);
CREATE INDEX IF NOT EXISTS idx_tasks_due ON tasks(due_date);
CREATE INDEX IF NOT EXISTS idx_points_child ON points_log(child_id);
