// 认证与视图导航：家长认证（Casdoor SSO）、学生/家长模式切换
import {
  parentToken, canAdmin, view, adminTab, pendingTab, authMode,
  currentStudent, students
} from "../state.js";
import { api } from "./useApi.js";
import { load, loadWrong } from "./useData.js";
import { loadAdmin } from "./useAdmin.js";

// 家长→学生：直接切换（免认证），并立即失效家长中心权限
export function backToStudent() {
  canAdmin.value = false; sessionStorage.removeItem("sp_admin"); view.value = "home";
}

// 学生→家长：必须重新认证。Casdoor 模式直接跳转 SSO，PIN 模式弹出登录框。
export async function enterAdmin(tab = "overview") {
  if (!canAdmin.value) { pendingTab.value = tab; location.href = "/api/parent/casdoor/login"; return; }
  adminTab.value = tab; view.value = "admin";
  await loadAdmin();
}

export function switchStudent(id) {
  currentStudent.value = id;
  localStorage.setItem("sp_stu", String(id));
  view.value = "home";
  return Promise.all([load(), loadWrong()]);
}

async function loadAuth() {
  try { authMode.value = (await api("/parent/auth-mode", { child: false })).mode || "casdoor"; } catch {}
}

// Casdoor 回跳后前端仅需恢复管理员身份；JWT 已持久化，loadAdmin 的 401 会自动清理
export function markAdmin() {
  canAdmin.value = true;
  sessionStorage.setItem("sp_admin", "1");
}

export function logout() {
  parentToken.value = ""; localStorage.removeItem("sp_parent_jwt");
  canAdmin.value = false; sessionStorage.removeItem("sp_admin");
  view.value = "home";
}

export { loadAuth };
