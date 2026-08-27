-- 000003 多邻国式练习场次与评分制度
-- 说明：SQLite 方言；切 Postgres/MySQL 时复制到 migrations/postgres/ 调整类型即可。

-- 练习场次：一次关卡 = 一个 session，记录科目/关卡/题目与结算
CREATE TABLE IF NOT EXISTS practice_sessions (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  child_id    INTEGER NOT NULL DEFAULT 1,
  subject     TEXT NOT NULL,             -- words | reading | math
  level       INTEGER NOT NULL DEFAULT 1,
  total       INTEGER NOT NULL DEFAULT 0,  -- 本关题目总数
  correct     INTEGER NOT NULL DEFAULT 0,  -- 答对题数
  max_combo   INTEGER NOT NULL DEFAULT 0,  -- 最高连击
  bonus       INTEGER NOT NULL DEFAULT 0,  -- 结算奖分（连击+星级）
  stars       INTEGER NOT NULL DEFAULT 0,  -- 结算星级 0-3，未结算为 0 且 finished=0
  finished    INTEGER NOT NULL DEFAULT 0,
  created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  finished_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sessions_child ON practice_sessions(child_id);
CREATE INDEX IF NOT EXISTS idx_sessions_subject ON practice_sessions(child_id, subject);

-- 场次内每题作答流水（供重做同一关时防重复计分）
CREATE TABLE IF NOT EXISTS session_answers (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id INTEGER NOT NULL,
  ref_id     INTEGER NOT NULL,           -- word/readings/math_problems 的 id 或 reading question id
  correct    INTEGER NOT NULL DEFAULT 0,
  points     INTEGER NOT NULL DEFAULT 0, -- 该题基础得分（不含奖分）
  combo      INTEGER NOT NULL DEFAULT 0, -- 作答时的连击数
  answered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_session_answers ON session_answers(session_id);
