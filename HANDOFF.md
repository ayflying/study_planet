# 学霸星球 StudyPlanet · 项目交接

> 最后更新：2026-08-27 ｜ 版本 0.1.5 ｜ 交接给新会话直接对接用

## 一句话概览

给 11 岁五年级学生做的多邻国式学习工作台：Go GoFrame v2 服务端 + 内嵌单文件前端 + SQLite + Docker Compose 部署在 217 服务器。

## 基本信息

| 项 | 值 |
|---|---|
| 仓库 | https://github.com/ayflying/study_planet.git（main 分支） |
| 本地路径 | `D:\git\yunloli\studyplanet` |
| 服务端框架 | Go 1.25 + GoFrame v2.10.3（ghttp + gcfg） |
| 数据库 | SQLite（modernc.org/sqlite 纯 Go 驱动）+ golang-migrate 内嵌 SQL 迁移（go:embed） |
| 镜像 | `ghcr.io/ayflying/study_planet-client` + `ghcr.io/ayflying/study_planet-server`，各自双标签 `latest` + 版本号 |
| 部署 | 217 测试服务器 192.168.50.217:18180；client Nginx 对外、server Go API 内网，纯拉取模式 |
| 数据卷 | `db` → server 挂载 `/app/data` |
| 前端 | `client/index.html` 独立静态工作台；client Nginx 托管并将 `/api/*` 同源反代到 server |
| 鉴权 | 孩子端开放接口（student_id 参数）；家长端 JWT + Casdoor OIDC（配置在 .env） |

## 目录结构

```
studyplanet/
├── client/                    # 独立工作台：index.html + assets/logo.png + Nginx 反代配置
│   ├── index.html             # 保留学习航线、关卡、宠物、积分、连击、星级与结算玩法
│   ├── nginx.conf             # /api 同源反代到 server
│   └── Dockerfile
├── server/                    # 独立 Go API：main.go、go.mod、internal/、manifest/、Dockerfile
│   └── internal/              # config / db / handler / middleware / model / seed
├── VERSION                    # 当前版本号，每次更新 +1（当前 0.1.5）
├── logo.png                   # 品牌 Logo 源文件（与 docs/logo.png、client/assets/logo.png 同源）
├── .env                       # 真实隐私配置（gitignore，勿提交）
├── docker-compose.yml         # client + server 镜像与 SQLite 数据卷
├── .github/workflows/docker.yml # push main 构建并推送两份 GHCR 镜像
└── docs/logo.png
```

## 部署流程（标准闭环）

1. 本地改代码 → `cd server && go build ./... && go vet ./...` 验证；客户端检查 `client/index.html` 与 Nginx 配置
2. VERSION +1 → 中文 commit（只提交本次涉及文件）→ `git push origin main`
3. GitHub Actions 自动构建并推送 client/server 两份镜像（各 latest + VERSION 号双标签），约 3 分钟
4. 217 服务器：`docker compose pull && docker compose up -d`（compose 在服务器上，路径见 README）
5. 验证：`curl http://192.168.50.217:18180/api/health` 应返回 `"version":"<新版本>"`
6. 回滚：compose 里 tag 改成旧版本号（如 `0.1.2`）再 `up -d`

## 配置说明

- **Casdoor**：endpoint=https://oidc.luoe.cn，回调地址不用配置——`server/internal/handler/casdoor.go` 的 `RedirectURIOf(r)` 每次请求按 Nginx 透传的 Host / X-Forwarded-* 实时推导。新增域名时只需在 Casdoor 应用侧把回调地址加入白名单。
- **PIN 回退**：.env 里 Casdoor 三项留空即进入 PIN 登录模式（PARENT_PIN）。
- **端口**：217 上 18080 被占，统一用 18180。

## 已知问题与约定（重要）

- **列表页/首页类派生数据优先 24h 缓存；数据库是事实来源。**
- **数据库迁移**：启动时自动跑（幂等），新增表结构只加新 migration 文件，不改旧的。
- **版本一致性**：VERSION 文件 = 镜像 tag = health 返回的 version（CI 通过 ldflags 注入）。
- **不要提交**：.env、bin/、data/*.db。
- **品牌 Logo**：正式 logo 为根目录 `logo.png`（地球+书本）。三处同源副本：根目录源文件、`docs/logo.png`（README 用）、`client/assets/logo.png`（页面与 favicon）。改 logo 时三处一起换。
- gf v2.10.3 注意：`gctx` 包不存在（用 context.Background()）；`gcfg.SetPath` 已删；`s.Run()` 无返回值。
- 历史 bug 修复记录见 git log，典型：NULL 扫描、Compose 卷前缀名；工作台现由 client Nginx 托管，server 不再绑定 `/`、`/app` 或 `/assets`。

## v0.1.5 更新记录（2026-08-27）

- 项目按职责拆分为 `client/` 与 `server/`，根目录不再混放工作台页面与 Go 服务端源码
- 完整保留此前学习星球工作台界面与现有玩法：学习航线、单词/阅读/数学闯关、宠物陪伴、积分、连击、星级、结算均未改动
- client Nginx 对外托管页面，并同源反代 `/api/*` 到 server；外部地址仍是 `:18180`
- CI 改为构建 `study_planet-client` 与 `study_planet-server` 两份版本一致的镜像

## v0.1.4 更新记录（2026-08-27）

- 正式启用品牌 Logo：根目录上传的 `logo.png`（地球+书本）替换旧 SVG（学士帽星球）
- `logo.png` 内嵌进二进制（`internal/handler/assets/logo.png`），路由 `/assets/logo.png`（GET+HEAD），旧 `/assets/logo.svg` 移除
- 修复首页 logo 一直 404 的 bug：app.html 原来拼出的 `/api/assets/logo.svg` 服务端从未注册，现改为服务端注入的站点级绝对路径 `LOGO_URL → /assets/logo.png`
- 页面补充 favicon（同用 logo.png）

## v0.1.3 修复记录（2026-08-27）

- 首页空白：app.html renderTabs 引用未定义变量 `home` 导致启动脚本中断（已删死代码）
- 关卡卡片无色：CSS `.u-word/.u-read` 与 JS 生成的 `u-words/u-reading` 类名不匹配（已对齐）
- HEAD / 返回 404：补 HEAD 绑定

## 下一步候选（未开工）

- [ ] 按 GoFrame v2 官方规范重构项目结构（当前是轻量自组织结构，非 gf init 标准布局：无 internal/controller、internal/service、internal/dao 分层）
- [ ] 家长端 Web 管理页（当前只有 API，无管理界面）
- [ ] 学生账号切换 UI 优化、更多关卡内容
- [ ] SQLite → PostgreSQL/MySQL 升级评估（迁移模块已备好）
