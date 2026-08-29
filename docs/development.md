# 开发指南

## 环境要求

- Go ≥ 1.24（服务端）
- Node ≥ 22（客户端构建）
- Docker + Docker Compose（容器化运行）
- 必须：MySQL 8+（唯一支持数据库）

## 本地开发

### 服务端

```bash
cd server
go run .            # 默认 :8080，需 DB_DSN 指向 MySQL
```

首次启动自动：建库 → 迁移建表 → 写种子数据。删除 `data/studyplanet.db` 即可重置。

### 客户端（热更新）

```bash
cd client
npm install
npm run dev         # Vite 开发服务器，/api 代理到本地 :8080
```

### 客户端生产构建

```bash
npm --prefix client run build   # 产物 client/dist，服务端打包时内嵌
```

## 常用验证命令（改动后必跑）

```bash
cd server
go vet ./...        # 静态检查
go test ./...       # 单元测试 + 编译验证

# 客户端
npm --prefix client run build

# 容器编排校验
docker compose config
```

## 新增接口流程

1. `server/internal/logic/studyplanet/`（handlers.go / practice.go 等）新增 `func (s *sStudyPlanet) Xxx(r *ghttp.Request)`，用 `s.ok(r, data)` / `s.fail(r, code, msg)` 返回。
2. 在 server 目录执行 `gf gen service`，重新生成 `internal/service` 接口与注册文件。
3. `server/internal/controller/studyplanet/router.go` 注册路由（`logic.Xxx`，logic 变量来自 `service.StudyPlanet()`）：
   - 孩子端 → `/api` 公开组（记得处理 `student_id`）；
   - 家长端 → `/api/parent` 鉴权组。
4. 需要新模型 → `internal/model/model.go` 加结构体（`db` + `json` tag）。
5. 跑 `go vet` / `go test`，curl 验证后更新 `docs/api.md`。

## 新增数据库迁移

迁移脚本目录：`server/internal/db/migrations/mysql/`。

```bash
# 1. 新版本号 = 现有最大版本 + 1
server/internal/db/migrations/mysql/000004_your_change.up.sql
```

注意事项：

- MySQL 脚本会被逐语句执行，**一条语句写一行**，不要用存储过程/触发器多语句块。
- `settings.key` 写成 `` `key` ``（保留字）。
- 迁移器按版本号顺序执行、跳过已应用版本；库版本回退会拒绝启动。
- 本地重置验证：建议新建空库或手动 `DROP TABLE` 后重启。

## MySQL SQL 写法要点

| 场景 | 写法 |
|---|---|
| 占位符 | 一律 `?` |
| upsert | `ON DUPLICATE KEY UPDATE` |
| 自增 ID | `res.LastInsertId()` |
| 可空唯一列 | 空值写 `NULL` 不写 `''` |
| 保留字 | `key`、`rank` 等加反引号 |

## 环境变量

本地调试可在启动前 export（或参考 `.env.example` 建 `.env`）：

```bash
export DB_DRIVER=mysql
export DB_DSN='user:pass@tcp(host:3306)/db?charset=utf8mb4&parseTime=true&loc=Local'
export PARENT_PIN=1234
export JWT_SECRET='random-long-string-at-least-32-chars'
```

## 代码风格约定

- 注释、commit message 用中文；标识符/术语保留英文。
- commit 只包含本次涉及文件。
- 服务端错误信息直接返回中文文案（面向前端展示）。
- 时间统一 `2006-01-02 15:04:05` 格式字符串。
- 禁止 Unicode 转义（`\uXXXX`），使用正常可读字符。

## 代码组织规范（硬性要求）

**单文件不要太大，目标不超过 200 行。** 超过时应按职责拆分：

- 服务端：按业务领域拆 logic 文件（如 `studyplanet/` 下按模块分文件），controller/handler 按路由组拆分。
- 客户端：按组件、页面、组合式函数（composable）拆分，公共工具函数收敛到 utils。

**可复用的方法必须抽象出去，不要重复代码。** 判断标准：同一段逻辑出现第二次时就该抽成公共函数/组件，而不是复制粘贴。抽象位置约定：

- 服务端：跨模块复用放 `internal/` 对应支撑包（如 contentlib）；仅本模块用的放同包内私有函数。
- 客户端：UI 复用抽组件，逻辑复用抽 composable/工具函数，接口请求统一走封装的 `api()`。
- 常量、错误码、枚举只定义一处，禁止多处硬编码同一魔法值。

动机：方便修改（改一处即全局生效）与复用（新功能直接组装现有能力）。

## 项目内文档维护

接口/表结构/部署变更时，同步更新：

- `docs/api.md`（接口）
- `docs/database.md`（表结构/迁移）
- `docs/deployment.md`（部署/CI）
- 根目录 `VERSION`（发版时）
