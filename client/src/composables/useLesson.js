// 练习流程：开始关卡、错题巩固混入、判分（服务端判分，前端无答案）、连击特效、键盘操作
import { computed } from "vue";
import {
  view, lesson, questions, questionIndex, combo, sessionId, result, fx,
  activeStudent, hasStudent, notice
} from "../state.js";
import { api, fail } from "./useApi.js";
import { shuffle } from "../utils.js";
import { load, loadWrong, starCount, totalStars, maxCombo, units, wrongSubject } from "./useData.js";
import { enterAdmin } from "./useAuth.js";

// 特效播放：答对/连击/答错/三星（kind 决定样式与持续时间）
export function playFx(kind) {
  fx.value = kind;
  setTimeout(() => { if (fx.value === kind) fx.value = ""; }, kind === "star" ? 1600 : 900);
}

// 内容库题目 → 练习题对象（选择题/填空题通用）
function toQuestion(q, review = false) {
  return {
    id: q.id, type: "content", subject: q.subject,
    title: q.question, passage: q.passage || "", subtitle: q.topic || "",
    options: shuffle(q.options || []),
    qtype: q.qtype || "choice",
    fillOpen: (q.qtype === "fill") || !(q.options || []).length,
    fillInput: "", fillSel: undefined,
    fillHint: q.qtype === "fill" ? "输入答案（支持分数如 3/4），回车提交" : "输入答案，回车提交",
    review
  };
}

export async function startUnit(unit) {
  if (!hasStudent.value) { await enterAdmin("students"); return; }
  try {
    const picked = await api(`/content/pick?subject=${unit.kind}&grade=${activeStudent.value?.grade || 1}&limit=5`) || [];
    if (!picked.length) throw new Error(`「${unit.title}」的题库还是空的，等待资料导入`);
    let list = picked.map(q => toQuestion(q));
    list = await mixWrongQuestions(unit.kind, list);
    if (!list.length) throw new Error("这一关还没有题目");
    const s = await api("/sessions", { method: "POST", body: { subject: unit.kind, level: activeStudent.value?.grade || 1, total: list.length } });
    lesson.value = unit; questions.value = list; questionIndex.value = 0;
    combo.value = 0; sessionId.value = s.id; result.value = null;
    view.value = "lesson";
  } catch (e) { fail(e); }
}

// 错题巩固：内容库题目通过 /content/item 回取（服务端判分，前端无答案）
async function mixWrongQuestions(kind, list) {
  try {
    const wrongs = await api(`/wrong-questions?subject=${kind}`) || [];
    const wrongIds = wrongs.filter(w => !Number(w.resolved)).slice(0, 3).map(w => Number(w.ref_id));
    const idSet = new Set(list.map(q => q.id));
    const reviewQs = [];
    for (const wid of wrongIds) {
      if (idSet.has(wid)) continue;
      try { reviewQs.push(toQuestion(await api(`/content/item?id=${wid}`), true)); } catch {}
    }
    // 每答 2 题新题穿插 1 道错题，其余错题排在末尾
    const mixed = [];
    let wi = 0;
    list.forEach((q, i) => { mixed.push(q); if ((i + 1) % 2 === 0 && wi < reviewQs.length) mixed.push(reviewQs[wi++]); });
    while (wi < reviewQs.length) mixed.push(reviewQs[wi++]);
    return mixed;
  } catch { return list; }
}

export const activeQuestion = computed(() => questions.value[questionIndex.value]);
export const progress = computed(() => questions.value.length ? Math.round((questionIndex.value / questions.value.length) * 100) : 0);

// 选择题判分（与 submitFill 共用答案应用逻辑）
async function applyAnswer(q, answer) {
  const r = await api("/content/answer", { method: "POST", body: { id: q.id, answer, session_id: sessionId.value } });
  q.correct = r.correct ?? false; q.rightAnswer = r.answer || "";
  combo.value = q.correct ? (r.combo || combo.value + 1) : 0;
  if (q.correct) { q.xp = r.xp || 0; playFx(combo.value >= 3 ? "combo" : "right"); } else { playFx("wrong"); }
}

export async function selectOption(o) {
  const q = activeQuestion.value;
  if (!q || q.picked !== undefined) return;
  q.picked = o;
  try { await applyAnswer(q, o); } catch (e) { fail(e); }
}

// 填空题输入（服务端判分，支持数值/分数/文本）
export async function submitFill(val) {
  const q = activeQuestion.value;
  if (!q || q.picked !== undefined || !val) return;
  q.picked = val;
  try { await applyAnswer(q, val); } catch (e) { fail(e); }
}

export async function nextQuestion() {
  if (questionIndex.value + 1 < questions.value.length) { questionIndex.value++; return; }
  try {
    result.value = await api(`/sessions/${sessionId.value}/finish`, { method: "POST" });
    if (result.value?.stars >= 2) playFx("star");
    if (result.value?.snack_name) notice.value = `🎁 获得零食：${result.value.snack_name} ×1`;
    await load(); await loadWrong();
  } catch (e) { fail(e); }
}

export function backHome() { result.value = null; view.value = "home"; }

// ---------- 全局键盘：练习页回车提交/下一题 ----------
export function lessonKeyHandler(e) {
  if (view.value !== "lesson" || e.key !== "Enter") return;
  const q = activeQuestion.value;
  if (!q) return;
  // 1) 填空题：输入框回车 = 提交答案
  if (q.fillOpen && q.picked === undefined && q.fillInput?.trim()) { e.preventDefault(); submitFill(q.fillInput.trim()); return; }
  // 2) 已判分：回车 = 下一题
  if (q.picked !== undefined) { e.preventDefault(); nextQuestion(); return; }
  // 3) 选择题：数字键选中后回车 = 确认
  if (q.options && q.fillSel !== undefined && !q.fillOpen) { e.preventDefault(); selectOption(q.options[q.fillSel]); }
}

export function lessonNumHandler(e) {
  if (view.value !== "lesson") return;
  const q = activeQuestion.value;
  if (!q || q.picked !== undefined || !q.options) return;
  if (/^[1-9]$/.test(e.key)) { const i = Number(e.key) - 1; if (q.options[i]) { q.fillSel = i; e.preventDefault(); } }
}

export { units, wrongSubject, starCount, totalStars, maxCombo };
