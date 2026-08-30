// api() 封装：所有接口请求的唯一出口（含 401 自动清理家长态、阻塞式 Loading 计数）
import { parentToken, loading, busy, error, notice, canAdmin, view, currentStudent } from "../state.js";
import { API_BASE } from "../constants.js";

export function fail(e) { error.value = e.message || String(e); }
export function clearError() { error.value = ""; }
export function clearNotice() { notice.value = ""; }

export function api(path, options = {}) {
  const sep = path.includes("?") ? "&" : "?";
  const headers = { "Content-Type": "application/json", ...(options.headers || {}) };
  if (parentToken.value) headers.Authorization = `Bearer ${parentToken.value}`;
  // 首次页面加载用原文本提示；其余请求走阻塞式 Loading 弹框（网慢时挡住交互防乱点）
  const counted = !loading.value;
  if (counted) busy.value++;
  return fetch(`${API_BASE}${path}${options.child === false ? "" : `${sep}student_id=${currentStudent.value}`}`, {
    ...options, headers, body: options.body ? JSON.stringify(options.body) : undefined
  }).finally(() => { if (counted) busy.value--; }).then(async r => {
    const data = await r.json().catch(() => null);
    // 401：JWT 失效，清理家长态回到学生模式（Casdoor 模式下重新进入家长中心会再走 SSO）
    if (r.status === 401 && parentToken.value) {
      parentToken.value = ""; localStorage.removeItem("sp_parent_jwt");
      canAdmin.value = false; sessionStorage.removeItem("sp_admin");
      view.value = "home";
    }
    // GF 标准响应：{code, message, data}，code !== 0 视为失败
    if (data && typeof data.code === "number") {
      if (data.code !== 0) throw new Error(data.message || `请求失败（${r.status}）`);
      return data.data;
    }
    if (!r.ok) throw new Error(data?.message || data?.error || `请求失败（${r.status}）`);
    return data;
  });
}
