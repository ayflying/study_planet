// 本地开发登录引导生成器：
// 用与本地服务端一致的 JWT_SECRET 自签一个长期 token（parent_id=1，10 年有效期），
// 并把 client/dist/local-login.html 写出来。浏览器访问该页即完成家长登录态注入。
// 仅用于本机开发环境（Casdoor 未配置时），不属于产品功能。
import { createHmac } from "node:crypto";
import { mkdirSync, writeFileSync, existsSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const SECRET = "local-dev-secret-0123456789abcdef";
const PARENT_ID = 1;
const TEN_YEARS_SECONDS = 10 * 365 * 24 * 3600;

const b64 = (obj) => Buffer.from(JSON.stringify(obj)).toString("base64url");
const header = b64({ alg: "HS256", typ: "JWT" });
const payload = b64({
  parent_id: PARENT_ID,
  exp: Math.floor(Date.now() / 1000) + TEN_YEARS_SECONDS,
});
const signature = createHmac("sha256", SECRET)
  .update(`${header}.${payload}`)
  .digest("base64url");
const jwt = `${header}.${payload}.${signature}`;

const repoRoot = join(dirname(fileURLToPath(import.meta.url)), "..");
const distDir = join(repoRoot, "client", "dist");
const html = `<!DOCTYPE html>
<html lang="zh-CN"><head><meta charset="UTF-8"><title>本地登录引导</title></head>
<body><p>正在写入本地登录态…</p>
<script>
try{
  localStorage.setItem('sp_parent_jwt','${jwt}');
  localStorage.setItem('sp_parent_name','本地家长');
}catch(e){alert('写入失败:'+e)}
location.href='/';
</script></body></html>
`;

try {
  mkdirSync(distDir, { recursive: true });
  writeFileSync(join(distDir, "local-login.html"), html, "utf8");
  if (!existsSync(join(distDir, "index.html"))) {
    console.warn("警告: client/dist/index.html 不存在，请先在 client/ 下执行 npm run build");
  }
  console.log("local-login.html 已生成:", join(distDir, "local-login.html"));
} catch (err) {
  console.error("生成失败:", err.message);
  process.exit(1);
}
