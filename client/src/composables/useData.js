// 全局数据加载：学生列表、学科、积分、练习场次、错题数、排行榜
import { computed } from "vue";
import {
  students, currentStudent, points, sessions, subjects,
  activeStudent, hasStudent, loading,
  leaderboard, wrongCount, showBoard
} from "../state.js";
import { api, fail } from "./useApi.js";
import { unitColors, fallbackUnits, wrongSubjectOrder, gradeLabel } from "../constants.js";

export const totalStars = () => sessions.value.reduce((s, x) => s + (x.stars || 0), 0);
export const maxCombo = () => sessions.value.reduce((b, x) => Math.max(b, x.max_combo || 0), 0);
export const starCount = (u) => sessions.value.filter(x => x.subject === u.kind && x.level === 1).reduce((b, x) => Math.max(b, x.stars || 0), 0);

// 首页数据加载：学生 → 学科 → 积分/场次
export async function load() {
  loading.value = true;
  try {
    const ss = await api("/students", { child: false });
    students.value = ss || [];
    if (!students.value.some(s => s.id === currentStudent.value)) currentStudent.value = students.value[0]?.id || 0;
    try {
      subjects.value = await api(`/subjects?grade=${activeStudent.value?.grade || 0}`, { child: false });
    } catch { subjects.value = []; }
    if (hasStudent.value) {
      const [pp, se] = await Promise.all([api("/points"), api("/sessions")]);
      points.value = pp?.total || 0;
      sessions.value = se || [];
    }
  } catch (e) { fail(e); } finally { loading.value = false; }
}

export async function loadWrong() {
  try {
    const w = await api("/wrong-questions");
    wrongCount.value = (w || []).filter(x => !Number(x.resolved)).length;
  } catch { wrongCount.value = 0; }
}

export async function loadLeaderboard() {
  try { leaderboard.value = await api("/leaderboard/weekly?limit=20"); } catch { leaderboard.value = null; }
}

export async function showLeaderboard() {
  await loadLeaderboard();
  leaderboard.value = leaderboard.value || { entries: [] };
  showBoard.value = true;
}

// 学习地图：学科从接口动态获取（内容库驱动，新增学科无需改前端）
export const units = computed(() => {
  const g = activeStudent.value?.grade || 0;
  // 只显示该学段开设的学科（后端 subjects 带 min/max_grade，这里再兜底过滤）
  const list = (subjects.value.length ? subjects.value : fallbackUnits).filter(s => !g || !s.min_grade || (g >= s.min_grade && g <= s.max_grade));
  return list.map((s, i) => ({ kind: s.code, title: s.name, sub: `${s.count || 0} 道题 · ${gradeLabel(activeStudent.value?.grade)}`, className: unitColors[i % unitColors.length], icon: s.icon || s.name[0] }));
});

export function wrongSubject() {
  for (const s of wrongSubjectOrder) {
    const u = subjects.value.find(x => x.code === s);
    if (u && u.count > 0) return s;
  }
  return "math";
}
