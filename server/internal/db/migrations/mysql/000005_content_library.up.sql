-- 000005 动态内容库：学科目录 + 统一题库（小学到初中全科）
-- 学习内容全部入库：以后新增/更新题目只导入数据，不改源码。

CREATE TABLE IF NOT EXISTS subjects (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(32) NOT NULL UNIQUE,
  name VARCHAR(64) NOT NULL,
  icon VARCHAR(16) NOT NULL DEFAULT '',
  color VARCHAR(32) NOT NULL DEFAULT '',
  min_grade INT NOT NULL DEFAULT 1,
  max_grade INT NOT NULL DEFAULT 9,
  sort INT NOT NULL DEFAULT 0,
  enabled TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS questions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  subject VARCHAR(32) NOT NULL,
  grade INT NOT NULL DEFAULT 1,
  topic VARCHAR(100) NOT NULL DEFAULT '',
  qtype VARCHAR(16) NOT NULL DEFAULT 'choice',
  passage TEXT NULL,
  question TEXT NOT NULL,
  options TEXT NOT NULL,
  answer VARCHAR(255) NOT NULL,
  explanation TEXT NULL,
  difficulty INT NOT NULL DEFAULT 1,
  source VARCHAR(64) NOT NULL DEFAULT '',
  content_hash CHAR(32) NOT NULL,
  enabled TINYINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_question_hash (content_hash),
  KEY idx_questions_pick (subject, enabled, grade),
  KEY idx_questions_grade (grade)
);
