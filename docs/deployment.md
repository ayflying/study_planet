# 部署运维

## 架构

单容器 `studyplanet-server`（GoFrame 托管 Vue 静态资源 + API），对外端口 `18180:8080`。数据库推荐外部 MySQL，SQLite 仅本地开发。

## Docker Compose

`docker-compose.yml`：

```yaml
services:
  server:
    image: ghcr.io/ayflying/study_planet:latest
    container_name: studyplanet-server
    restart: unless-stopped
    env_file: [.env]
    environment: [TZ=Asia/Shanghai]
    ports: ["18180:8080"]
    volumes: [db:/app/data]     # SQLite 模式数据卷；MySQL 模式不依赖
    healthcheck:
      test: ["CMD","wget","-qO-","http://localhost:8080/api/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
```

本地构建运行：

```bash
docker compose up -d --build
curl http://localhost:18180/api/health
```

## 环境变量（.env）

从 `.env.example` 复制。关键项：

| 变量 | 必填 | 说明 |
|---|---|---|
| `DB_DRIVER` | 是 | `mysql`（生产） / `sqlite`（开发） |
| `DB_DSN` | 是 | MySQL DSN 或 SQLite 文件路径 |
| `PARENT_PIN` | 建议 | 初始 PIN（仅空库种子写入；之后用接口改） |
| `JWT_SECRET` | 必须 | ≥32 位随机串；**更换后所有登录态失效** |
| `CASDOOR_ENDPOINT/CLIENT_ID/CLIENT_SECRET` | 可选 | 三项齐备切换 Casdoor 登录 |
| `CASDOOR_ORG_NAME / CASDOOR_APP_NAME` | 可选 | 默认 `built-in` / `studyplanet` |

## CI/CD（GitHub Actions）

`.github/workflows/docker.yml`：

- 触发：push 到 `main`（或手动 `workflow_dispatch`）。
- 流程：Node 构建 Vue → Go 构建服务端（内嵌 `client/dist`）→ 推送镜像 `ghcr.io/ayflying/study_planet`，`latest` + `VERSION` 版本号双标签。
- 发版：改根目录 `VERSION` → commit → push → CI 自动构建。

## 217 测试服务器

入口：`http://192.168.50.217:18180/`（内网测试环境，与正式环境严格区分）。

更新流程：

```bash
cd /opt/studyplanet
docker compose pull
docker compose up -d
docker compose ps          # 确认 healthy
curl http://localhost:18180/api/health
```

注意事项：

- `.env` 存放真实 MySQL 连接与 `JWT_SECRET`，**不入 Git**。
- 迁移在容器启动时自动执行：空库全量建表 + 种子；已有库按 `schema_migrations` 增量升级；版本回退会启动失败（查看 `docker compose logs`）。
- 临时 Cloudflare Quick Tunnel 已删除，不要重建；正式域名由用户自行配置（反代到 18180 即可，回调地址自动识别 `X-Forwarded-*`）。

## 健康检查与排障

| 检查 | 命令/位置 |
|---|---|
| 容器状态 | `docker compose ps`（应 `healthy`） |
| 健康接口 | `curl http://<host>:18180/api/health` |
| 启动日志 | `docker compose logs -f server` |
| 迁移失败 | 日志出现 `migrate` 报错；常见原因：MySQL 不可达、DSN 错、版本回退 |
| 登录全部失效 | `JWT_SECRET` 被更换 |

## 回滚

```bash
# 镜像标签改为历史版本号（如 0.1.9）
# 编辑 docker-compose.yml 中 image tag 后：
docker compose up -d
```

数据库回滚需谨慎：迁移器检测到库版本高于脚本版本会拒绝启动；确需回滚先手工把 `schema_migrations` 版本降级并核对表结构。

## 发布检查清单

- [ ] `go vet` / `go test` 通过
- [ ] `npm --prefix client run build` 通过
- [ ] `docker compose config` 通过
- [ ] `VERSION` 已更新
- [ ] 真实环境验证：health 200、登录可用、核心链路（创建学生→练习→积分）正常
- [ ] 容器 `healthy`，日志无迁移报错
