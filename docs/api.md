# 接口文档

Base URL：`http://<host>:18180/api`

## 通用约定

- 数据格式均为 JSON；错误统一返回 `{"error": "原因"}` + 对应 HTTP 状态码。
- 孩子端接口统一用查询参数 `student_id` 指定学生，**缺省为 1**。
- 家长接口需请求头 `Authorization: Bearer <token>`；token 通过登录接口获取，长期有效（不设过期），仅主动退出/更换 `JWT_SECRET` 失效。
- CORS 全开放（`*`），便于本地开发调试。

## 认证与状态

### GET /health

健康检查。响应：`{"status":"ok","version":"0.1.9"}`

### GET /parent/auth-mode

返回登录模式。响应：

```json
{"mode": "pin"}       // 未配置 Casdoor
{"mode": "casdoor"}   // 三项 Casdoor 配置齐备
```

### POST /parent/login（仅 PIN 模式）

请求：`{"pin": "1234"}`

响应：`{"token": "<jwt>"}`；PIN 错误 401；Casdoor 模式下 400。

### GET /parent/casdoor/login

302 跳转 Casdoor 授权页（回调地址按实际访问 host 自动推导，支持 `X-Forwarded-*` 反代头）。

### GET /parent/casdoor/callback

Casdoor 回调：授权码换用户信息 → upsert `parents` 表 → 302 回前端并携带 token（前端写入 `localStorage.sp_parent_jwt`）。

## 学生

### GET /students

学生列表（孩子选身份/家长切换用）。

```json
[
  {"id":1,"name":"小朋友","username":"","avatar":"🚀","grade":5,"created_at":"2026-08-28 10:00:00"}
]
```

### POST /parent/students 🔒

请求：`{"name":"二宝","username":"kid2","avatar":"🐯","grade":3}`（username 可空）

响应：创建后的学生对象。非空 username 重复 → 500。

### PUT /parent/students/:id 🔒

字段可选更新：`{"name":"...","username":"...","avatar":"...","grade":3}`

### DELETE /parent/students/:id 🔒

删除学生及其全部学习数据；**至少保留一个学生**，否则 400。

## 学习内容（孩子端）

### GET /words?level=

单词列表：`[{"id":1,"level":5,"word":"because","meaning":"因为","phonetic":"/bɪˈkɒz/","example":"..."}]`

### GET /words/:id?student_id=

单词详情，附该学生的掌握状态 `{"known":true,"last_reviewed":"..."}`。

### POST /words/:id/progress?student_id=

请求：`{"known":true,"session_id":0}`

- `known=true`：掌握 +5 分；若传 `session_id>0` 则走练习场次计分（连击+奖分）。
- 响应：`{"correct":true,"combo":3,"base_points":7,"combo_bonus":2}`

### GET /readings/:id

短文 + 理解题：

```json
{"id":1,"title":"龟兔赛跑","content":"...","level":5,
 "questions":[{"id":1,"question":"...","option_a":"...","option_b":"...","option_c":"...","option_d":"..."}]}
```

### POST /readings/:id/answer?student_id=

请求：`{"question_id":1,"answer":"它骄傲睡觉","session_id":0}`

响应：`{"correct":true,"points":2}`（答对 +2；支持 `answer` 大小写/空格容错）。

### GET /math?level=

数学题列表（含选项与解析字段，`answer` 不下发前端）。

### POST /math/:id/answer?student_id=

请求：`{"answer":"12","session_id":0}`；答对 +3，返回对错与正确答案/解析。

## 动态内容库（全科题库，内容入库不改源码）

学习内容统一存于 `questions` 表。前端学习地图从 `/subjects` 动态渲染；出题走内容库；判分走统一接口，复用连击/XP/错题本链路。

### GET /subjects

学科目录 + 每科题量：

```json
[{"code":"math","name":"数学","icon":"∑","color":"#27ae60","min_grade":1,"max_grade":9,"count":298},
 {"code":"physics","name":"物理","icon":"⚛","count":10}]
```

内置学科：english / chinese / math（1-9 年级）、physics（8-9）、chemistry（9）、biology / history / geography（7-9）。

### GET /content/pick?subject=math&grade=5&limit=5

随机抽题（`answer` 不下发前端）。`grade` 为学生年级，实际取 `grade±1` 范围保证题量；`limit` 1~20。

```json
[{"id":115,"subject":"math","grade":5,"topic":"小数运算","qtype":"choice",
  "question":"计算：1.29 + 3.24 = ?","options":["4.53","5.53","5.03","4.93"],"difficulty":1}]
```

### GET /content/item?id=

按 id 回取单题（不含答案），错题本巩固复习用。

### POST /content/answer?student_id=

统一判分。请求：`{"id":115,"answer":"4.53","session_id":9}`

- `session_id>0`：走场次计分（连击阶梯、XP 1:1、错题登记/消除），响应含 `combo/xp/review`。
- `session_id=0`：独立判分，答对 +3 分，响应 `{"correct":true,"answer":"4.53","explanation":"..."}`。

### POST /parent/content/import 🔒

**通用题目导入**——以后采集新学习资料只调此接口，不改源码。单次 ≤2000 题，按 `content_hash`（subject+question+answer 的 MD5）去重：

```json
{"questions":[{
  "subject":"math","grade":7,"topic":"有理数","qtype":"choice",
  "passage":"",                    // 可选：阅读短文
  "question":"计算：(-3) + (-5) = ?",
  "options":["-8","8","-2","2"],
  "answer":"-8",
  "explanation":"同号相加",         // 可选
  "difficulty":1,                  // 可选，默认 1
  "source":"自定义来源"             // 可选
}]}
```

响应：`{"imported":1,"skipped":1,"total":2}`（skipped 为重复跳过）。

### GET /parent/content/stats 🔒

内容库统计：`{"total":447,"subjects":[{"code":"math",...,"count":299},...]}`

内置题库：服务启动时若 `questions` 表为空，自动导入 446 道内置全科题（数学程序化生成保证答案正确、英语分级词表、语文古诗词/成语/文学常识、理科基础概念）。之后以数据库为准，重复启动不覆盖。

## 练习场次（闯关）

### POST /sessions?student_id=

开启一关。请求：`{"subject":"words","level":5,"total":5}`（total 1~50，缺省 5）

响应：场次对象 `{"id":9,"child_id":1,"subject":"words","level":5,"total":5,"correct":0,"max_combo":0,"bonus":0,"stars":0,"finished":0,...}`

### GET /sessions?student_id=

最近练习记录（成长统计用）。

### POST /sessions/:id/finish?student_id=

结算。星级规则：正确率 ≥90% 三星 / ≥70% 两星 / ≥50% 一星；三星额外 +10 分、两星 +5 分；同一关只结算一次（重复请求 400）。

### 连击规则（recordAnswer）

- 连击 = 最近一次答错后的连续正确数。
- 阶梯奖分：连击达 3/5/8/10 时一次性 +2/+4/+6/+8。
- 作答必须属于该学生且场次未结束，否则 400。

## 任务与积分

### GET /tasks?student_id=&status=

任务列表；逾期未完成实时计算为 `overdue` 状态（前端标红）。

### POST /tasks/:id/complete?student_id=

完成任务，奖励 `task.points` 积分；已完成任务重复提交 400。

### GET /points?student_id=

`{"total":42,"today_earned":10,"student_id":1}`

### GET /points/log?student_id=

积分流水：`[{"id":7,"delta":5,"reason":"单词认读:+5","created_at":"..."}]`

## 奖励与兑换

### GET /rewards

`[{"id":1,"name":"看动画 30 分钟","cost_points":20,"status":"active"}]`

### POST /rewards/:id/redeem?student_id=

提交兑换 → 生成 `redemptions(pending)`；积分不足 400。

### POST /parent/rewards 🔒

新增奖励：`{"name":"...","cost_points":50}`

### POST /parent/redemptions/:id/confirm 🔒

家长确认兑换：状态改 `confirmed` 并**扣除学生积分**。

## 任务管理（家长）

### POST /parent/tasks 🔒

`{"title":"练字一页","type":"daily","due_date":"2026-08-30","points":10,"student_id":1}`

### DELETE /parent/tasks/:id 🔒

删除任务。

## 家长安全

### POST /parent/set-pin 🔒

修改 PIN：`{"pin":"5678"}`；bcrypt 写 `settings.parent_pin`。

## 静态资源

| 路径 | 说明 |
|---|---|
| `/` | Vue 单页应用（内嵌构建产物） |
| `/assets/...` | 前端 JS/CSS/图片 |
| `/assets/logo.png` | 品牌 Logo |

## 状态码汇总

| 码 | 场景 |
|---|---|
| 400 | 参数错误 / Casdoor 模式下调 PIN 登录 / 重复结算 / 积分不足 |
| 401 | PIN 错误 / JWT 缺失或无效（前端收到 401 自动清理会话） |
| 404 | 学生/资源不存在 |
| 500 | 服务端错误（含数据库约束冲突） |

🔒 = 需家长 JWT。
