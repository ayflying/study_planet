# 系统架构

## 总体形态

单容器应用：GoFrame 服务端同时提供 API 与 Vue 静态资源，不使用 Nginx。

```text
浏览器 / 手机
    │  http://<host>:18180
    ▼
┌─────────────────────────────┐
│  studyplanet-server 容器     │
│  ┌───────────────────────┐  │
│  │ GoFrame (Go)          │  │
│  │  /            → Vue 静态资源（client/dist，内嵌 packed）
│  │  /assets/...  → 前端构建产物
│  │  /api/...     → REST API
│  └───────────────────────┘  │
└──────────────┬──────────────┘
               │ sqlx + ? 占位符
               ▼
     MySQL（推荐，远程库）或 SQLite（默认，容器卷 /app/data）
```

## 技术选型

| 层 | 技术 | 说明 |
|---|---|---|
| 客户端 | Vue 3 + Vite | 单文件主组件 `App.vue` + 全局 `style.css`，无路由库 |
| 服务端 | Go + GoFrame v2 | `gf init` 脚手架分层：cmd / controller / service / model |
| 数据访问 | jmoiron/sqlx | 业务查询用 `?` 占位符，SQLite/MySQL 通用 |
| 数据库 | MySQL（推荐）/ SQLite | 驱动：`go-sql-driver/mysql`、`modernc.org/sqlite`（CGO-free） |
| 迁移 | 自研轻量迁移器 | SQL 内嵌 `embed.FS`，启动时逐语句执行，版本表 `schema_migrations` |
| 认证 | Casdoor SSO（OIDC）+ PIN 回退 | 服务端签发 JWT（HS256），长期有效直至主动退出 |
| 部署 | Docker Compose + GitHub Actions | 镜像 `ghcr.io/ayflying/study_planet`，`latest` + 版本双标签 |

## 目录结构

```text
studyplanet/
├── client/                          # Vue 3 + Vite 客户端
│   ├── src/App.vue                  # 主组件：家长端 + 学生端全部页面
│   ├── src/main.js                  # Vue 入口
│   ├── src/style.css                # 设计系统 + 响应式样式
│   ├── public/logo.png              # 品牌 Logo
│   ├── vite.config.js               # 开发环境 /api 代理
│   └── package.json
├── server/                          # GoFrame 服务端
│   ├── main.go                      # 入口，调用 gcmd 命令
│   ├── internal/
│   │   ├── cmd/cmd.go               # 服务初始化：配置→数据库→迁移→种子→路由→静态资源
│   │   ├── config/config.go         # yaml + 环境变量配置加载（env 优先）
│   │   ├── controller/studyplanet/router.go  # 路由注册（唯一路由清单）
│   │   ├── service/studyplanet/
│   │   │   ├── handlers.go          # 学生/任务/积分/奖励/单词/阅读/数学 业务
│   │   │   ├── practice.go          # 练习场次：连击、星级、XP 结算
│   │   │   └── casdoor.go           # Casdoor OIDC 登录
│   │   ├── middleware/auth.go       # CORS + 家长 JWT 鉴权
│   │   ├── model/model.go           # 数据模型（db/json tag）
│   │   ├── db/
│   │   │   ├── db.go                # Open（双驱动）+ Migrate（自研迁移器）
│   │   │   └── migrations/
│   │   │       ├── sqlite/          # 000001~000003（SQLite 方言）
│   │   │       └── mysql/           # 000001（MySQL 全量建表）
│   │   └── seed/seed.go             # 空库种子数据（五年级示例）
│   ├── manifest/config/config.yaml  # 默认配置
│   ├── examples/gf-init-example/    # gf init 官方脚手架参考示例
│   └── Dockerfile                   # 两阶段：Node 构建 Vue → Go 构建 → 运行镜像
├── docs/                            # 项目文档（本目录）
├── docker-compose.yml               # 单容器编排，端口 18180:8080
├── .env.example                     # 环境变量模板（.env 不入库）
├── .github/workflows/docker.yml     # CI：push main 自动构建推送镜像
├── VERSION                          # 当前版本号（如 0.1.9）
└── Makefile                         # 常用命令封装
```

## 请求链路

```text
HTTP 请求
  → GoFrame Server
    → /api 分组中间件 CORS
      → 公开路由（学生自选身份、登录、学习内容）
      → /api/parent 分组中间件 ParentAuth（校验 Bearer JWT）
        → 家长管理接口
  → controller（router.go 绑定）
    → service/studyplanet（业务 + SQL）
      → sqlx → MySQL / SQLite
```

## 启动流程（internal/cmd/cmd.go）

1. `config.Load()`：读 yaml，环境变量覆盖。
2. `db.Open(driver, dsn)`：按驱动建立连接，Ping 验证。
3. `db.Migrate`：确保 `schema_migrations` 表 → 比对版本 → 逐语句执行缺失迁移（MySQL 按语句拆分，遇注释/空行跳过）。
4. `seed.Run`：`children` 表为空时写入示例数据（幂等）。
5. `BindRoutes`：注册 API 路由。
6. 绑定静态资源：Vue 构建产物（`internal/packed` 内嵌）+ `/assets`。
7. 启动 HTTP 监听 `SERVER_PORT`（默认 8080）。

## 关键设计决策

- **家长先行**：未登录只见家长登录入口；登录后才能创建学生，学生切换后进入学习工作台。
- **多学生隔离**：所有学习数据挂 `child_id`，孩子端请求统一带 `?student_id=`（缺省 1）。
- **登录长期有效**：JWT 不设 `exp`，前端 `localStorage` 持久保存；仅主动退出、清理站点数据或更换 `JWT_SECRET` 才会登出；服务端 401 时前端自动清理失效会话。
- **双数据库支持**：迁移脚本按方言分目录，业务 SQL 用通用 `?` 占位符，切换成本低。
- **迁移失败即启动失败**：版本回退或 SQL 错误直接报错退出，不静默跳过，避免破坏数据。
