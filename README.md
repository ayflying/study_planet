# 🪐 学霸星球 StudyPlanet

> 给小学生的学习闯关台服务端：单词卡片 · 语文阅读 · 数学题目 · 每日任务 · 积分奖励，家长用 Casdoor 登录，一个家长可创建并切换多个学生账号。

![logo](docs/logo.svg)

## 功能总览

- **六大学习模块**：单词卡片（翻转记忆+掌握加分）、语文阅读（短文+理解题判分）、数学题目（即答即讲）、每日任务（逾期自动标红、完成加分）、积分（累计/今日）、奖励与兑换（提交兑换→家长确认扣分）。
- **多学生账号**：家长可创建多个学生（姓名/用户名/头像/年级），孩子端通过 `?student_id=` 切换身份，任务、积分、单词掌握进度、兑换记录全部按学生隔离；至少保留一个学生。
- **家长认证**：
  - **Casdoor SSO**（推荐）：配置环境变量后自动启用，OIDC 授权码流程，家长信息落库 `parents` 表，签发本站 JWT；
  - **PIN 回退**：未配置 Casdoor 时自动启用 PIN 登录（bcrypt 存储）。
- **迁移模块**：`golang-migrate` + 内嵌 SQL，启动自动升级，版本表 `schema_migrations`；SQLite 现用，日后可平滑切 Postgres/MySQL。

## 快速开始

```bash
# 本地运行
go run .          # 默认 :8080，SQLite ./data/studyplanet.db

# Docker Compose
docker compose up -d --build
curl http://localhost:18180/api/health
```

默认家长 PIN：`1234`（仅 PIN 模式；务必在部署时用 `PARENT_PIN` 环境变量修改）。

## API 一览

### 公开接口

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/health` | 健康检查 |
| GET | `/api/assets/logo.svg` | 品牌 Logo |
| GET | `/api/students` | 学生列表（孩子选身份） |
| GET | `/api/parent/auth-mode` | 登录模式 `pin` / `casdoor` |
| POST | `/api/parent/login` | PIN 登录 → `{token}`（仅 PIN 模式） |
| GET | `/api/parent/casdoor/login` | 302 跳 Casdoor 授权页 |
| GET | `/api/parent/casdoor/callback` | 授权码换 JWT（浏览器回调） |
| GET | `/api/words?level=` | 单词列表 |
| GET | `/api/words/:id?student_id=` | 单词详情（含该学生掌握状态） |
| POST | `/api/words/:id/progress?student_id=` | 标记掌握 `{known:true}` +5 分 |
| GET | `/api/readings/:id` | 阅读详情+理解题 |
| POST | `/api/readings/:id/answer?student_id=` | 提交答案 `{question_id,answer}` 对+2 分 |
| GET | `/api/math?level=` | 数学题列表 |
| POST | `/api/math/:id/answer?student_id=` | 提交答案 `{answer}` 对+3 分 |
| GET | `/api/tasks?student_id=&status=` | 任务列表（逾期计算返回 `overdue`） |
| POST | `/api/tasks/:id/complete?student_id=` | 完成任务加分 |
| GET | `/api/points?student_id=` | 积分 `{total,today_earned,student_id}` |
| GET | `/api/points/log?student_id=` | 积分流水 |
| GET | `/api/rewards` | 奖励列表 |
| POST | `/api/rewards/:id/redeem?student_id=` | 提交兑换（积分足额才受理） |

> `student_id` 缺省为 1，兼容旧客户端。

### 家长接口（需 `Authorization: Bearer <token>`）

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/parent/tasks` | 发布任务 `{title,type,due_date,points,student_id}` |
| DELETE | `/api/parent/tasks/:id` | 删除任务 |
| POST | `/api/parent/rewards` | 新增奖励 `{name,cost_points}` |
| POST | `/api/parent/redemptions/:id/confirm` | 确认兑换（扣对应学生积分） |
| POST | `/api/parent/set-pin` | 修改家长 PIN `{pin}` |
| POST | `/api/parent/students` | 创建学生 `{name,username?,avatar?,grade?}` |
| PUT | `/api/parent/students/:id` | 修改学生（字段可选更新） |
| DELETE | `/api/parent/students/:id` | 删除学生（至少保留一个；清空其数据） |

## Casdoor 接入

在 Casdoor 控制台创建应用（Redirect URL 填 `http://<host>:<port>/api/parent/casdoor/callback`），然后给容器加环境变量：

```yaml
environment:
  - CASDOOR_ENDPOINT=https://casdoor.example.com
  - CASDOOR_CLIENT_ID=xxxxxxxx
  - CASDOOR_CLIENT_SECRET=xxxxxxxx
  - CASDOOR_ORG_NAME=built-in        # 可选
  - CASDOOR_APP_NAME=studyplanet     # 可选，默认 studyplanet
  - CASDOOR_REDIRECT_URL=http://<your-host>:<port>/api/parent/casdoor/callback  # 可选，缺省按请求推导
```

三项核心（endpoint/client id/secret）配齐 → `GET /api/parent/auth-mode` 返回 `"mode":"casdoor"`，前端把「PIN 输入框」换成「跳转 Casdoor 登录」按钮即可；PIN 登录同时被禁用（400）。回调成功后本站 JWT 写入 `localStorage.sp_parent_jwt`。

家长首次 SSO 登录自动落库 `parents(casdoor_sub, display_name, avatar, last_login_at)`。

## 配置

配置文件 `manifest/config/config.yaml`，同名环境变量优先：

| 环境变量 | 默认 | 说明 |
|---|---|---|
| `SERVER_PORT` | 8080 | 监听端口 |
| `DB_DSN` | data/studyplanet.db | SQLite 文件路径 |
| `PARENT_PIN` | 1234 | 种子 PIN（仅首库初始化写库） |
| `JWT_SECRET` | change-me-in-prod | 务必改随机长字符串 |
| `CASDOOR_*` | 空 | 见上节 |

## 数据库与未来切库

迁移文件内嵌于二进制（`internal/db/migrations/*.sql`），启动幂等执行。新增迁移：

```bash
cp internal/db/migrations/000002_multi_student.up.sql   internal/db/migrations/000003_your_change.up.sql   # 改内容
cp internal/db/migrations/000002_multi_student.down.sql internal/db/migrations/000003_your_change.down.sql
```

当前 SQLite 用纯 Go 驱动 `modernc.org/sqlite`（CGO-free）。切换 Postgres/MySQL：

1. `go get` 并 import 对应驱动；
2. `db.Open` / `Migrate` 换驱动名与 DSN；
3. 将 `migrations/` 复制为 `migrations/postgres/`，调整自增（`AUTOINCREMENT`→`BIGSERIAL`）与时间戳类型。

业务查询走 sqlx 占位符 `?`，主流库通用。

## 发版与 CI

1. 改 `VERSION` 文件（如 `0.1.2` → `0.1.3`），提交并 push 到 `main`；
2. GitHub Actions 自动构建镜像，推送到 GHCR 双标签：`latest` + 版本号（如 `0.1.3`）；
3. 服务器上 `docker compose pull && docker compose up -d` 即更新到 latest。

回滚到指定版本：

```bash
# 修改 docker-compose.yml 的 image 标签为历史版本号
sed -i 's|study_planet:latest|study_planet:0.1.1|' docker-compose.yml
docker compose up -d
```

CI 触发方式：push 到 `main` 自动触发；也可在 Actions 页面手动 `workflow_dispatch`。

## 目录结构

```
studyplanet/
├── main.go                      # 入口与路由注册
├── internal/
│   ├── config/config.go         # yaml + env 配置加载
│   ├── db/
│   │   ├── db.go                # 连接 & golang-migrate 启动迁移
│   │   └── migrations/*.sql     # 版本化 SQL（内嵌）
│   ├── handler/
│   │   ├── handlers.go          # 学习/任务/积分/兑换/学生 CRUD
│   │   ├── casdoor.go           # OIDC 授权码登录
│   │   └── assets.go            # Logo 内嵌服务
│   ├── middleware/auth.go       # CORS + 家长 JWT 鉴权
│   ├── model/model.go           # 数据模型
│   └── seed/seed.go             # 空库种子（示例内容）
├── manifest/config/config.yaml
├── docs/logo.svg                # 品牌 Logo 源文件
├── Dockerfile                   # 多阶段构建 CGO_ENABLED=0
└── docker-compose.yml           # 卷持久化 + healthcheck
```

## 在线实例

部署后访问 `http://<your-host>:<port>/api/health` 验证；宿主端口冲突时改 `docker-compose.yml` 的 `ports` 即可；SQLite 落卷 `db`，重建容器不丢数据。
