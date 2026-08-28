# 数据库文档

## 连接配置

配置优先级：环境变量 > `server/manifest/config/config.yaml`。

| 环境变量 | 默认 | 说明 |
|---|---|---|
| `DB_DRIVER` | `sqlite` | `sqlite` 或 `mysql` |
| `DB_DSN` | `data/studyplanet.db` | SQLite 为文件路径；MySQL 为 DSN |

MySQL DSN 示例（`.env`，不入库）：

```text
DB_DRIVER=mysql
DB_DSN=用户名:密码@tcp(主机:3306)/库名?charset=utf8mb4&parseTime=true&loc=Local
```

生产当前使用 MySQL（远程库），容器部署通过 `.env` 注入；SQLite 仅用于本地快速开发。

## 迁移机制（自研轻量迁移器）

实现在 `server/internal/db/db.go`，启动时自动执行：

1. 确保版本表 `schema_migrations(version INT PRIMARY KEY)` 存在。
2. 读当前库版本，与内嵌迁移脚本（`embed.FS`）比对。
3. 空库 → 全量建表；落后 → 按版本号逐个增量执行。
4. **版本回退**（库版本 > 脚本版本）→ 报错拒绝启动，不静默破坏数据。
5. 每个版本脚本按语句拆分逐条执行（MySQL 驱动不支持多语句），已执行版本跳过。

迁移目录按方言组织：

```text
server/internal/db/migrations/
├── sqlite/            # 000001_init / 000002_multi_student / 000003_practice_sessions
└── mysql/             # 000001_init（全量建表，对应 SQLite 三个版本合并后的最终结构）
```

新增迁移步骤见 [development.md](development.md)。

> 历史备注：曾尝试 golang-migrate v4.19.1，其 MySQL 驱动不支持多语句且失败后留下 dirty 标记导致迁移被跳过，故替换为逐语句自研迁移器。

## 表结构（MySQL，共 15 张表）

### children 学生档案

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK AUTO_INCREMENT | |
| name | VARCHAR(255) NOT NULL DEFAULT '小朋友' | 显示名 |
| username | VARCHAR(191) NULL | 登录用户名，可空；非空唯一（唯一索引，NULL 不参与约束） |
| avatar | VARCHAR(32) NOT NULL DEFAULT '🚀' | emoji 头像 |
| grade | INT NOT NULL DEFAULT 5 | 年级 |
| created_at | TIMESTAMP DEFAULT CURRENT_TIMESTAMP | |

查询输出用 `COALESCE(username,'') AS username` 归一为字符串。

### words 单词

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| level | INT DEFAULT 1 | 难度/年级（索引 idx_words_level） |
| word | VARCHAR(255) NOT NULL | 单词 |
| meaning | TEXT NOT NULL | 中文释义 |
| phonetic | VARCHAR(255) DEFAULT '' | 音标 |
| example | TEXT NOT NULL | 例句 |
| created_at | TIMESTAMP | |

### word_progress 单词掌握进度

| 列 | 类型 | 说明 |
|---|---|---|
| word_id + child_id | 复合主键 | 每学生每单词一行 |
| known | TINYINT DEFAULT 0 | 是否掌握 |
| last_reviewed | TIMESTAMP NULL | 最近复习时间 |

### readings 阅读短文

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| title | VARCHAR(255) NOT NULL | 标题 |
| content | TEXT NOT NULL | 正文 |
| level | INT DEFAULT 1 | 难度 |

### reading_questions 阅读理解题

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| reading_id | BIGINT（索引） | 所属短文 |
| question | TEXT NOT NULL | 题干 |
| option_a/b/c/d | TEXT NOT NULL | 四个选项 |
| answer | VARCHAR(255) NOT NULL | 正确答案（选项文本） |

### math_problems 数学题

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| level | INT DEFAULT 1（索引） | 难度 |
| type | VARCHAR(100) NOT NULL | 题型 |
| question | TEXT NOT NULL | 题干 |
| options | TEXT NOT NULL | 选项（JSON/分隔文本） |
| answer | VARCHAR(255) NOT NULL | 正确答案 |
| explanation | TEXT NOT NULL | 解析 |

### tasks 任务

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| title | VARCHAR(255) NOT NULL | 任务名 |
| type | VARCHAR(100) NOT NULL | 类型 |
| due_date | DATE NULL | 截止日（过期未完成为 `overdue`，由接口实时计算） |
| points | INT DEFAULT 0 | 完成奖励积分 |
| status | VARCHAR(32) DEFAULT 'pending' | pending / done |
| child_id | BIGINT DEFAULT 1（索引） | 所属学生 |
| created_at / completed_at | TIMESTAMP | |

### points_log 积分流水

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| child_id | BIGINT（索引） | 学生 |
| delta | INT | 变动值（正负） |
| reason | VARCHAR(255) | 事由（如 `单词认读:+5`） |
| created_at | TIMESTAMP | 今日积分按 `created_at` 当日聚合 |

### rewards 奖励

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| name | VARCHAR(255) NOT NULL | 奖励名 |
| cost_points | INT DEFAULT 0 | 兑换所需积分 |
| status | VARCHAR(32) DEFAULT 'active' | active / off |

### redemptions 兑换记录

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| reward_id | BIGINT | 兑换的奖励 |
| child_id | BIGINT DEFAULT 1（索引） | 学生 |
| status | VARCHAR(32) DEFAULT 'pending' | pending / confirmed |
| requested_at | TIMESTAMP | 申请时间 |
| confirmed_at | TIMESTAMP NULL | 家长确认时间（确认时扣分） |

### settings 键值配置

| 列 | 类型 | 说明 |
|---|---|---|
| key | VARCHAR(255) PK | 注意：MySQL 保留字，SQL 中需写反引号 `` `key` `` |
| value | TEXT NOT NULL | 目前存 `parent_pin`（bcrypt 哈希） |

### parents 家长账号（Casdoor）

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| casdoor_sub | VARCHAR(255) UNIQUE | Casdoor 用户唯一标识 |
| display_name / avatar | VARCHAR | 昵称/头像 |
| created_at | TIMESTAMP | |
| last_login_at | TIMESTAMP NULL | 每次登录更新 |

### practice_sessions 练习场次（闯关核心）

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| child_id | BIGINT（索引） | 学生 |
| subject | VARCHAR(32) | words / reading / math |
| level | INT DEFAULT 1 | 关卡难度 |
| total | INT DEFAULT 0 | 本关题数 |
| correct | INT DEFAULT 0 | 答对数 |
| max_combo | INT DEFAULT 0 | 最高连击 |
| bonus | INT DEFAULT 0 | 连击阶梯奖分累计 |
| stars | INT DEFAULT 0 | 结算星级 0~3 |
| finished | TINYINT DEFAULT 0 | 是否已结算 |
| created_at / finished_at | TIMESTAMP | |

### session_answers 场次作答明细

| 列 | 类型 | 说明 |
|---|---|---|
| id | BIGINT PK | |
| session_id | BIGINT（索引） | 所属场次 |
| ref_id | BIGINT | 题目 ID（单词/问题/数学题） |
| correct | TINYINT | 是否答对 |
| points | INT | 本题得分（含连击加成） |
| combo | INT | 作答后连击数 |
| answered_at | TIMESTAMP | |

## 种子数据（空库自动写入，幂等）

`children` 有数据则跳过。写入内容：

- 默认学生「小朋友」（头像 🚀，五年级）
- 10 个五年级单词（because、beautiful、favorite 等，含音标例句）
- 2 篇寓言阅读（龟兔赛跑、守株待兔）各带 2 道理解题
- 4 道数学题
- 3 个任务（含 1 个逾期示例）
- 3 个奖励
- 家长 PIN：取 `PARENT_PIN`（默认 `1234`），bcrypt 后写 `settings.parent_pin`

## 计分规则（与表相关）

| 行为 | 积分 | 落表 |
|---|---|---|
| 单词标记掌握 | +5 | word_progress / points_log |
| 阅读理解题答对 | +2 | points_log |
| 数学题答对 | +3 | points_log |
| 连击阶梯 3/5/8/10 | +2/+4/+6/+8 | session_answers.points |
| 关卡结算三星/两星 | +10/+5 | practice_sessions.stars |

## MySQL 兼容性要点（写 SQL 时注意）

1. `settings.key` 是保留字 → 一律写 `` `key` ``。
2. `children.username` 空值写 `NULL` 不写 `''`（否则触发唯一冲突）。
3. upsert 用 `ON DUPLICATE KEY UPDATE`（SQLite 为 `ON CONFLICT ... DO UPDATE`）。
4. 取自增 ID 用 `res.LastInsertId()`。
5. 时间统一 `2006-01-02 15:04:05` 格式字符串 + `parseTime=true`。
