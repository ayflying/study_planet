# 学霸星球 StudyPlanet 项目文档

> 给小学生的学习闯关工作台：单词卡片 · 语文阅读 · 数学题目 · 每日任务 · 积分奖励。家长登录后创建/切换多个学生账号，孩子在学习地图上闯关拿星。

## 文档目录

| 文档 | 内容 |
|---|---|
| [architecture.md](architecture.md) | 系统架构：技术选型、目录结构、前后端协作、请求链路 |
| [database.md](database.md) | 数据库文档：连接配置、迁移机制、全部表结构、种子数据 |
| [api.md](api.md) | 接口文档：全部 HTTP 接口的路径、参数、响应示例、错误格式 |
| [design.md](design.md) | 界面设计：设计语言、页面结构、主流程、响应式规则 |
| [development.md](development.md) | 开发指南：本地开发、新增接口/迁移、测试与验证 |
| [deployment.md](deployment.md) | 部署运维：Docker Compose、CI/CD、217 服务器、回滚 |
| [roadmap.md](roadmap.md) | 开发计划：已完成里程碑、进行中、后续规划 |

## 快速入口

- 本地运行：`cd server && go run .`（默认 `:8080`，需 `DB_DSN` 指向 MySQL）
- Docker 运行：`docker compose up -d --build`，访问 `http://localhost:18180/`
- 健康检查：`GET /api/health`
- 当前版本：见根目录 `VERSION` 文件
