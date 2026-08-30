// 家长中心：数据加载、任务/奖励发布、孩子账号管理
import {
  canAdmin, parentToken, tasks, rewards, logs, notice,
  newTask, newReward, newStudent, showAddStudent,
  showMyTasks, myTasks, newStudentTask, showMyRewards, myRewards,
  students, currentStudent, busyRedeemSnack, points, starTotal
} from "../state.js";
import { api, fail } from "./useApi.js";
import { load } from "./useData.js";

export async function loadAdmin() {
  if (!canAdmin.value || !parentToken.value) return;
  try {
    const [tt, rr, ll] = await Promise.all([
      api(`/tasks?student_id=${currentStudent.value}`, { child: false }),
      api("/rewards", { child: false }),
      api(`/points/log?student_id=${currentStudent.value}`, { child: false })
    ]);
    tasks.value = tt || []; rewards.value = rr || []; logs.value = ll || [];
  } catch (e) { fail(e); }
}

export async function addTask() {
  try {
    await api("/parent/tasks", { method: "POST", child: false, body: { ...newTask.value, student_id: currentStudent.value } });
    newTask.value = { title: "", type: "学习", due_date: "", points: 5 };
    notice.value = "任务已发布";
    await loadAdmin();
  } catch (e) { fail(e); }
}

export async function addReward() {
  try {
    await api("/parent/rewards", { method: "POST", child: false, body: newReward.value });
    newReward.value = { name: "", cost_points: 20 };
    notice.value = "奖励已添加";
    await loadAdmin();
  } catch (e) { fail(e); }
}

export async function addStudent() {
  if (!newStudent.value.name.trim()) { fail(new Error("请先填写孩子姓名")); return; }
  try {
    await api("/parent/students", { method: "POST", child: false, body: newStudent.value });
    showAddStudent.value = false;
    newStudent.value = { name: "", username: "", avatar: "🐣", grade: 1 };
    notice.value = "孩子账号已创建，可以在底部「学习」页切换到 TA 的学习工作台";
    await load(); await loadAdmin();
  } catch (e) { fail(e); }
}

export async function deleteStudent(id) {
  if (!confirm("确定删除这个学生及其学习数据吗？")) return;
  try {
    await api(`/parent/students/${id}`, { method: "DELETE", child: false });
    await load(); await loadAdmin();
  } catch (e) { fail(e); }
}

export async function deleteTask(id) {
  try {
    await api(`/parent/tasks/${id}`, { method: "DELETE", child: false });
    await loadAdmin();
  } catch (e) { fail(e); }
}

export function openAddStudent() {
  newStudent.value = { name: "", username: "", avatar: "🐣", grade: 1 };
  showAddStudent.value = true;
}

// ---------- 学生自己的功能（无需家长认证）：完成任务 / 兑换奖励 ----------
export async function openMyTasks() {
  try { myTasks.value = await api("/tasks") || []; showMyTasks.value = true; } catch (e) { fail(e); }
}

export async function openMyRewards() {
  try {
    // 先刷新积分和星星
    const pp = await api("/points");
    if (pp?.total !== undefined) points.value = pp.total;
    if (pp?.star_total !== undefined) starTotal.value = pp.star_total;
    myRewards.value = await api("/rewards") || [];
    showMyRewards.value = true;
  } catch (e) { fail(e); }
}

export async function doStudentAddTask() {
  try {
    await api("/tasks", { method: "POST", body: { ...newStudentTask.value, student_id: currentStudent.value } });
    newStudentTask.value = { title: "", type: "学习", due_date: "", points: 5 };
    notice.value = "任务已创建 ✍️";
    await openMyTasks();
  } catch (e) { fail(e); }
}

export async function doStudentDeleteTask(id) {
  try {
    await api(`/tasks/${id}`, { method: "DELETE", body: { student_id: currentStudent.value } });
    notice.value = "任务已删除";
    await openMyTasks();
  } catch (e) { fail(e); }
}

export async function doCompleteTask(id) {
  try {
    const r = await api(`/tasks/${id}/complete`, { method: "POST" });
    notice.value = "任务完成，积分已到账 🎉" + (r?.snack_name ? `，掉落零食：${r.snack_name} ×1` : "");
    await openMyTasks(); await load();
  } catch (e) { fail(e); }
}

export async function doRedeem(id) {
  try {
    const r = await api(`/rewards/${id}/redeem`, { method: "POST" });
    notice.value = r.message || "已提交兑换，等待家长确认";
    await openMyRewards();
  } catch (e) { fail(e); }
}

export async function doRedeemSnack(stars) {
  try {
    busyRedeemSnack.value = true;
    const r = await api("/redeem-snack", { method: "POST", body: { stars } });
    notice.value = `🎁 兑换成功！获得零食：${r.snack_name} ×1，剩余 ${r.stars_left} 颗星星`;
    await load();
  } catch (e) { fail(e); } finally { busyRedeemSnack.value = false; }
}
