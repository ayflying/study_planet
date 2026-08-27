# 学霸星球 StudyPlanet · 项目交接

> 最后更新：2026-08-27 ｜ 版本 0.1.3 ｜ 交接给新会话直接对接用

## 一句话概览

给 11 岁五年级学生做的多邻国式学习工作台：Go GoFrame v2 服务端 + 内嵌单文件前端 + SQLite + Docker Compose 部署在 217 服务器。

## 基本信息

| 项 | 值 |
|---|---|
| 仓库 | https://github.com/ayflying/study_planet.git（main 分支） |
| 本地路径 | `D:\git\yunloli\studyplanet` |
| 服务端框架 | Go 1.25 + GoFrame v2.10.3（ghttp + gcfg） |
| 数据库 | SQLite（modernc.org/sqlite 纯 Go 驱动）+ golang-migrate 内嵌 SQL 迁移（go:embed） |
| 镜像 | ghcr.io/ayflying/study_planet，双标签 `latest` + 版本号（如 `0.1.3`） |
| 部署 | 217 服务器 192.168.50.217:18180（容器内 8080），纯拉取模式（compose 无 build） |
| 数据卷 | `studyplanet_db` → 挂载 `/app/data` |
| 前端 | `internal/handler/assets/app.html` 单文件（go:embed 打进二进制），`GET /` 和 `GET /app` 提供 |
| 鉴权 | 孩子端开放接口（student_id 参数）；家长端 JWT + Casdoor OIDC（配置在 .env） |

## 目录结构

```
studyplanet/
├── main.go                    # 入口：路由注册（/api 组 + 站点级 GET/HEAD /、/app、/assets/logo.svg）
├── VERSION                    # 当前版本号，每次更新 +1（当前 0.1.3）
├── .env                       # 真实隐私配置（gitignore，勿提交）：CASDOOR_ENDPOINT/CLIENT_ID/SECRET 等
├── .env.example               # 模板（入库）
├── docker-compose.yml         # 仅 image: ghcr.io/ayflying/study_planet:latest + env_file
├── Dockerfile                 # golang:1.25-alpine 多阶段，ARG APP_VERSION 注入 ldflags
├── .github/workflows/docker.yml  # push main 触发 GHCR 双标签构建
├── internal/
│   ├── config/config.go       # env 读取：SERVER_PORT / DB_DSN / PARENT_PIN / JWT_SECRET / CASDOOR_*
│   ├── db/                    # Open + Migrate；migrations/000001~000003（init、multi_student、practice_sessions）
│   ├── handler/               # handlers.go / practice.go / casdoor.go / assets.go + assets/app.html、logo.svg
│   ├── middleware/auth.go     # CORS + ParentAuth(JWT)
│   ├── model/model.go         # 数据模型
│   └── seed/seed.go           # 空库写入五年级示例数据（单词 10 / 阅读 / 数学 4 题）
└── docs/logo.svg
```

## 部署流程（标准闭环）

1. 本地改代码 → `go build ./... && go vet ./...` 验证
2. VERSION +1 → 中文 commit（只提交本次涉及文件）→ `git push origin main`
3. GitHub Actions 自动构建镜像推 GHCR（latest + VERSION 号双标签），约 3 分钟
4. 217 服务器：`docker compose pull && docker compose up -d`（compose 在服务器上，路径见 README）
5. 验证：`curl http://192.168.50.217:18180/api/health` 应返回 `"version":"<新版本>"`
6. 回滚：compose 里 tag 改成旧版本号（如 `0.1.2`）再 `up -d`

## 配置说明

- **Casdoor**：endpoint=https://oidc.luoe.cn，回调地址不用配置——`internal/handler/casdoor.go` 的 `RedirectURIOf(r)` 每次请求按 Host / X-Forwarded-* 实时推导。新增域名时只需在 Casdoor 应用侧把回调地址加入白名单。
- **PIN 回退**：.env 里 Casdoor 三项留空即进入 PIN 登录模式（PARENT_PIN）。
- **端口**：217 上 18080 被占，统一用 18180。

## 已知问题与约定（重要）

- **列表页/首页类派生数据优先 24h 缓存；数据库是事实来源。**
- **数据库迁移**：启动时自动跑（幂等），新增表结构只加新 migration 文件，不改旧的。
- **版本一致性**：VERSION 文件 = 镜像 tag = health 返回的 version（CI 通过 ldflags 注入）。
- **不要提交**：.env、bin/、data/*.db、logo.png（未跟踪，本地草稿）。
- gf v2.10.3 注意：`gctx` 包不存在（用 context.Background()）；`gcfg.SetPath` 已删；`s.Run()` 无返回值。
- 历史 bug 修复记录见 git log，典型：NULL 扫描、Compose 卷前缀名、HEAD 404（已绑 HEAD:/ 与 HEAD:/app）。

## v0.1.3 修复记录（2026-08-27）

- 首页空白：app.html renderTabs 引用未定义变量 `home` 导致启动脚本中断（已删死代码）
- 关卡卡片无色：CSS `.u-word/.u-read` 与 JS 生成的 `u-words/u-reading` 类名不匹配（已对齐）
- HEAD / 返回 404：补 HEAD 绑定

## 下一步候选（未开工）

- [ ] 按 GoFrame v2 官方规范重构项目结构（当前是轻量自组织结构，非 gf init 标准布局：无 internal/controller、internal/service、internal/dao 分层）
- [ ] 家长端 Web 管理页（当前只有 API，无管理界面）
- [ ] 学生账号切换 UI 优化、更多关卡内容
- [ ] SQLite → PostgreSQL/MySQL 升级评估（迁移模块已备好）
