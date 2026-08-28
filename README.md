# 🪐 学霸星球 StudyPlanet

> 给小学生的学习闯关工作台：客户端与服务端独立部署。单词卡片 · 语文阅读 · 数学题目 · 每日任务 · 积分奖励，家长用 Casdoor 登录，一个家长可创建并切换多个学生账号。

![logo](docs/logo.png)

## 项目分层

- `client/`：Vue 3 + Vite 学习星球工作台，按组件维护页面和玩法；构建产物由 GoFrame 静态文件服务直接托管。
- `server/`：GoFrame API、静态资源服务、SQLite 数据库迁移、种子数据与 Casdoor 登录。
- 外部访问地址维持不变：`http://<host>:18180/`，单容器同时提供页面与 API。

## 功能总览

- **六大学习模块**：单词卡片（翻转记忆+掌握加分）、语文阅读（短文+理解题判分）、数学题目（即答即讲）、每日任务（逾期自动标红、完成加分）、积分（累计/今日）、奖励与兑换（提交兑换→家长确认扣分）。
- **多学生账号**：家长可创建多个学生（姓名/用户名/头像/年级），孩子端通过 `?student_id=` 切换身份，任务、积分、单词掌握进度、兑换记录全部按学生隔离；至少保留一个学生。
- **家长认证**：
  - **Casdoor SSO**（推荐）：配置环境变量后自动启用，OIDC 授权码流程，家长信息落库 `parents` 表，签发本站 JWT；
  - **PIN 回退**：未配置 Casdoor 时自动启用 PIN 登录（bcrypt 存储）。
- **迁移模块**：`golang-migrate` + 内嵌 SQL，启动自动升级，版本表 `schema_migrations`；SQLite 现用，日后可平滑切 Postgres/MySQL。

## 快速开始

```bash
# 本地运行服务端 API（默认 :8080，SQLite ./data/studyplanet.db）
cd server && go run .

# Docker Compose：GoFrame server 同时托管 Vue 静态资源与 API
docker compose up -d --build
curl http://localhost:18180/api/health
```

默认家长 PIN：`1234`（仅 PIN 模式；务必在部署时用 `PARENT_PIN` 环境变量修改）。

## API 一览

完整接口文档（参数、响应示例、错误码）见 [docs/api.md](docs/api.md)。

### 公开接口

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/health` | 健康检查 |
| GET | `/assets/logo.png` | 品牌 Logo（PNG） |
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
```

回调地址**无需配置**：授权与回调阶段都会按用户实际访问的地址自动生成 `http(s)://<访问host>/api/parent/casdoor/callback`，并识别 `X-Forwarded-Proto` / `X-Forwarded-Host` 反代头，换域名/端口/走反向代理都不用改配置。

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

## 数据库

支持 MySQL（生产推荐）与 SQLite（本地开发）。迁移机制、全部表结构、种子数据与 MySQL 兼容要点见 [docs/database.md](docs/database.md)。

新增迁移：在 `server/internal/db/migrations/sqlite/` 与 `migrations/mysql/` 同步添加下一版本号脚本，启动自动执行。开发细节见 [docs/development.md](docs/development.md)。

## 发版与 CI

1. 改 `VERSION` 文件，提交并 push 到 `main`；
2. GitHub Actions 自动构建统一镜像 `study_planet`（GoFrame API + Vue 静态资源），推送 `latest` + 版本号双标签；
3. 服务器上 `docker compose pull && docker compose up -d` 更新单容器应用。

回滚到指定版本：

```bash
# 将 study_planet image 标签改为历史版本号，例如 0.1.6
docker compose up -d
```

CI 触发方式：push 到 `main` 自动触发；也可在 Actions 页面手动 `workflow_dispatch`。

## 目录结构

```
studyplanet/
├── client/                       # 独立 Vue 3 + Vite 客户端源码
│   ├── src/App.vue               # 学习航线、闯关、宠物陪伴、积分、连击、结算组件
│   ├── src/main.js / style.css   # Vue 入口与样式
│   ├── assets/logo.png           # 客户端品牌 Logo
│   ├── vite.config.js            # 开发环境 /api 代理
│   └── package.json              # 前端依赖与构建命令
├── server/                       # GoFrame 服务端：API + 静态资源托管
│   ├── main.go                   # API 路由注册
│   ├── internal/
│   │   ├── config/config.go      # yaml + env 配置加载
│   │   ├── db/                   # 连接与内嵌迁移 SQL
│   │   ├── cmd/                  # GoFrame 命令入口
│   │   ├── controller/studyplanet/ # API 路由
│   │   ├── service/studyplanet/  # 学习、家长、Casdoor 业务
│   │   ├── middleware/           # CORS + 家长 JWT 鉴权
│   │   ├── db/ migrate/          # 数据库连接与迁移
│   │   └── seed/                 # 空库示例数据
│   ├── examples/gf-init-example/ # gf init 生成的官方脚手架示例
│   ├── manifest/config/          # 服务端配置
│   ├── go.mod / go.sum
│   └── Dockerfile                # Vue 构建 + GoFrame 运行镜像
├── docs/                         # 项目文档（架构/数据库/接口/设计/开发/部署/计划）
├── VERSION                        # 双镜像统一版本
└── docker-compose.yml            # 单 GoFrame 容器 + SQLite 数据卷
```

## 文档

完整项目文档在 [docs/](docs/README.md)：

| 文档 | 内容 |
|---|---|
| [架构](docs/architecture.md) | 技术选型、目录结构、请求链路、启动流程 |
| [数据库](docs/database.md) | 连接配置、迁移机制、15 张表结构、种子数据 |
| [接口](docs/api.md) | 全部 HTTP 接口参数与响应示例 |
| [界面设计](docs/design.md) | 设计语言、页面结构、主流程、响应式规则 |
| [开发指南](docs/development.md) | 本地开发、新增接口/迁移、双数据库兼容写法 |
| [部署运维](docs/deployment.md) | Docker Compose、CI/CD、217 服务器、回滚 |
| [开发计划](docs/roadmap.md) | 已完成里程碑与后续规划 |

## 在线实例

部署后访问 `http://<your-host>:<port>/api/health` 验证；宿主端口冲突时改 `docker-compose.yml` 的 `ports` 即可；SQLite 落卷 `db`，重建容器不丢数据。
