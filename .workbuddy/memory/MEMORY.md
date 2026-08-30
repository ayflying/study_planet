# study_planet 项目约定

## 版本与交付流程（用户明确约定，必须执行）
- **每次修改完内容 → `VERSION` 文件版本号 +1（如 0.2.1 → 0.2.2）→ git commit 并 push 到 origin**。
- **自动推送规则（2026-08-30 用户明确要求）**：每次写完代码改动后，自动执行 VERSION +1、git commit、git push 全流程，不再向用户确认"是否要推送"——直接推送并核实。只有遇到无法解决的网络错误时才告知用户。
- 推送网络坑：环境变量代理（127.0.0.1:62305）会 CONNECT 502，需 `env -u http_proxy -u https_proxy` 清除；GitHub 直连 IP 时通时断，用 `for ip in ...; do timeout 3 bash -c "echo > /dev/tcp/$ip/443" && break; done` 轮询找活 IP，再 `-c http.curloptResolve="github.com:443:<IP>"` 推送；amend 后 force-with-lease 需显式 `--force-with-lease=main:<SHA>`。
- **凭据坑**：push 会弹「CredentialHelperSelector」图形窗口卡死推送 → 加 `-c credential.helper= -c credential.helper=store` 绕过选择器直接用 ~/.git-credentials 已存凭据（liusihua123@github.com）。
- git 提交身份：repo-local 已配置 user.name=liusihua, user.email=1129985028@qq.com（历史提交曾是 ayflying <anyang@adesk.com>，已按用户要求改用 QQ 邮箱）。
- 推送命令完整模板：`env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY GIT_TERMINAL_PROMPT=0 timeout 240 git -c credential.helper= -c credential.helper=store -c http.curloptResolve="github.com:443:<活IP>" push origin main`
- **bash 假成功坑（2026-08-30 确认）**：Git Bash 的 git 网络操作（push/fetch/ls-remote）在网络差时会 exit 0 且零输出但实际没连上 GitHub。**git 网络操作一律用 PowerShell 执行并核实输出**（PowerShell 的 git 能给出真实错误），推送后必须 `ls-remote` 或 `fetch` 确认远程 SHA 变化才算成功。
- 远程节奏：远程（ayflying 侧）提交频繁且可能与我们并行开发，每次推送前必须重新 fetch+合并，不可盲目 push。
- **备用代理通道（用户 2026-08-30 提供，当前对 GitHub 不可用）**：`http://100.66.1.7:10808`——TCP 通但 CONNECT 隧道被 abort（git 报 "Proxy CONNECT aborted"）。直连 IP 轮询全失败时可重试 `git -c http.proxy=http://100.66.1.7:10808 push origin main`。
- **packed-refs 惯性坑**：本机 origin/main 引用（在 .git/packed-refs 里）经常不随 fetch/ls-remote 自动更新，导致本地状态误报 ahead/behind。推送后核对远程 SHA（ls-remote 真实输出）时，若本地引用不符，直接编辑 .git/packed-refs 校正。

## 技术栈
- 后端 GoFrame v2（Go 便携版 C:\Users\liusihua\.workbuddy\binaries\go\go\bin\go.exe，用户已自装 1.27；GOPROXY=https://goproxy.cn,direct）。
- 数据库：**仅 MySQL**（远程 2026-08-29 提交 0606165 已彻底移除 SQLite 支持：migrations/sqlite 目录删除、go.mod 去掉 glebarez 驱动、DB_DSN 必填）。本地测试也走 MySQL（.env 配 DB_DSN）。
- 前端 Vue3 + vite，已按规范拆分（v0.3.9）：App.vue 仅壳（51行），state.js 全局状态 / constants.js 常量 / utils.js 工具 / composables/（useApi/useAuth/useData/useLesson/useBattle/usePet/useAdmin）/ views/（Home/Lesson/Battle/BattleHome/Pet/Admin/GlobalModals）。改前端时先找对应 composable/视图，不要再往 App.vue 堆代码。
- 认证：Casdoor SSO（家长身份由 SSO 建立，PIN 登录已废弃）。

## 代码组织规范（用户明确要求，必须遵守）
- **禁止把代码全写到一个文件**：新功能先按模块划分落点再动手；单文件目标 ≤200 行，超了按职责拆分（后端按业务领域拆 logic 文件，前端按组件/composable/utils 拆分）。
- **可复用方法必须抽象出去**：同一段逻辑出现第二次就抽公共函数/组件，禁止复制粘贴；常量/错误码/魔法值只定义一处。
- 该规范已写入 docs/development.md「代码组织规范（硬性要求）」一节，含提交前自查清单。

## 宠物零食系统（v0.4.5）
- **有限库存**：投喂零食改为有限数量，初始全0；完成任务/对战胜利/学习闯关可掉落。
- **概率掉落**：权重法——不掉落40、苹果40、小鱼干30、牛奶30、星星糖30、蛋糕20、火锅10（总权重200），每把最多1个最少0个。
- **惩罚机制**：饱食度到0→清空积分（删points_log）；好感度到0→清空星星（sessions.stars/max_combo归零）。
- **库存列**：pets 表 food_apple/food_fish/food_milk/food_star/food_cake/food_hotpot（迁移 000008）。
- **星星兑换商店**（v0.4.8）：奖励商店新增"⭐ 零食兑换"区，用星星兑换宠物零食——5星换小鱼干/牛奶/星星糖（随机三选一）、8星换蛋糕、15星换小火锅。pets 表加 stars_spent 列记录已花费星星（迁移 000009）。
- **学习探索地图不再掉落零食**（v0.4.8）：移除随机掉落，改为奖励商店星星兑换。

## 进行中功能
- 填空题题库（contentgen/gen_fill.go）、真人对战（5题×10秒、速度计分每题10分、段位、结算奖励）、宠物模式（投喂/好感度/有限零食）、答题页回车键支持。
