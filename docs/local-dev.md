# 本地开发环境指南

> 不依赖 Docker / Casdoor / 内网，在本机把「学霸星球」完整跑起来。

## 一键启动

双击运行：

```
scripts\start-local.bat
```

它会自动：

1. 启动便携版 MySQL（端口 `33061`，已在运行则跳过）
2. 启动服务端（端口 `8095`，已在运行则跳过）
3. 打开浏览器进入登录引导页

停止一切：

```
scripts\stop-local.bat
```

## 首次使用（只需一次）

浏览器访问 <http://127.0.0.1:8095/local-login.html> —— 它会把本地家长登录态写入浏览器，然后自动跳回首页。

之后直接访问 <http://127.0.0.1:8095/> 即可。

## 前置条件（已完成，无需重复）

| 组件 | 位置 | 说明 |
|---|---|---|
| 便携版 MySQL 8.0.29 | `%USERPROFILE%\.workbuddy\binaries\mysql\mysql-8.0.29-winx64` | 端口 33061，root 空密码，业务账号 `sp` / `sp123456` |
| 服务端二进制 | `server\studyplanet_new.exe` | 由 `server` 源码编译 |
| 前端产物 | `client\dist` | 由 `client` 源码 `npm run build` 生成 |
| 登录引导页 | `client\dist\local-login.html` | 由 `scripts\prepare-local-login.mjs` 生成 |

## 重新构建（改了代码后）

```bash
# 前端
cd client && npm run build

# 服务端
cd server && go build -o studyplanet_new.exe .

# 重新生成登录引导页（JWT 有效期 10 年，一般不用重跑）
node scripts/prepare-local-login.mjs
```

## 本地登录说明（为什么需要引导页）

项目已停用 PIN 登录，正式环境走 Casdoor SSO；本机没有 Casdoor，
所以用与本地服务端一致的 `JWT_SECRET` 自签了一个 `parent_id=1` 的长期
token，由引导页写入 `localStorage`（与服务端 Casdoor 回调页同样的机制）。

数据库里预置了对应数据：家长 `local-dev`（id=1）、学生「小朋友」（5 年级）。

## 常见问题

**Q: 端口被占用？**
`8095` / `33061` 被占时改 `scripts\start-local.bat` 顶部的 `PORT` / `MYSQL_PORT`
（MySQL 端口还需同步改 my.ini）。

**Q: 白屏？**
服务端必须从 `server` 目录启动（脚本已处理）——它按工作目录找 `../client/dist`；
若从仓库根启动会命中 `client/index.html` 的开发版导致白屏。

**Q: 数据想清空重来？**
停止服务后删除 `%USERPROFILE%\.workbuddy\binaries\mysql\mysql-8.0.29-winx64\data`
再运行启动脚本，迁移会全量重建（种子题库自动导入）。
