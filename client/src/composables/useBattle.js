// 真人对战：WebSocket 状态机（匹配 → 抢答 → 结算）、段位页数据
import {
  view, battlePhase, battleMe, battleOpp, battleQs, battleIdx, battleRemain,
  battlePicked, battleCorrect, battleGain, battleResult, battleRank, battleHistory,
  currentStudent, activeStudent, hasStudent
} from "../state.js";
import { WS_URL } from "../constants.js";
import { api, fail } from "./useApi.js";

// WebSocket 实例不需要响应式，用普通变量持有即可
let battleWS = null;

export function battleSend(obj) { if (battleWS?.readyState === 1) battleWS.send(JSON.stringify(obj)); }

export function startBattle() {
  if (!hasStudent.value) return;
  battlePhase.value = "matching"; battleResult.value = null; battlePicked.value = undefined; battleCorrect.value = false;
  battleMe.value = { total: 0 }; battleOpp.value = { total: 0, name: "", avatar: "", isBot: false }; battleIdx.value = -1;
  const ws = new WebSocket(WS_URL);
  battleWS = ws;
  ws.onopen = () => ws.send(JSON.stringify({ type: "join", student_id: currentStudent.value, subject: "math", grade: activeStudent.value?.grade || 1 }));
  ws.onmessage = ev => { try { battleMsg(JSON.parse(ev.data)); } catch {} };
  ws.onclose = () => { if (battlePhase.value !== "finished") battleWS = null; };
  ws.onerror = () => { if (battlePhase.value === "matching") { fail(new Error("对战服务连接失败，请确认服务端已启动")); battlePhase.value = "idle"; } };
}

function battleMsg(m) {
  if (m.type === "matched") {
    battleOpp.value = m.opponent || {}; battleQs.value = m.questions || [];
  } else if (m.type === "question_next") {
    battleIdx.value = m.qindex; battleRemain.value = m.remain || 10; battlePicked.value = undefined; battleCorrect.value = false; battleGain.value = 0;
    battlePhase.value = "fighting";
  } else if (m.type === "tick") {
    if (m.qindex === battleIdx.value && m.remain !== undefined) battleRemain.value = m.remain;
    if (m.opp_total !== undefined) battleOpp.value = { ...battleOpp.value, total: m.opp_total };
  } else if (m.type === "opp_done") {
    if (m.opp_total !== undefined) battleOpp.value = { ...battleOpp.value, total: m.opp_total };
  } else if (m.type === "answer_result") {
    battlePicked.value = m.qindex; battleCorrect.value = m.correct; battleGain.value = m.score || 0;
    battleMe.value = { ...battleMe.value, total: m.total };
    if (m.opp_total !== undefined) battleOpp.value = { ...battleOpp.value, total: m.opp_total };
  } else if (m.type === "finished") {
    battleResult.value = m; battlePhase.value = "finished";
    try { battleWS?.close(); } catch {} battleWS = null;
  }
}

export function battleAnswer(ans) {
  if (battlePicked.value !== undefined || battlePhase.value !== "fighting") return;
  battleSend({ type: "answer", qindex: battleIdx.value, answer: ans });
}

export function battleNext() { /* 服务端自动推进下一题，无需操作 */ }

export function exitBattle() {
  try { battleWS?.close(); } catch {} battleWS = null;
  battlePhase.value = "idle";
}

// 对战选项键位（1-9 数字键选中；已判分回车进入下一题由服务端自动推进）
export function battleKeyHandler(e) {
  if (view.value !== "battle") return;
  const q = battleQs.value[battleIdx.value];
  if (battlePhase.value === "fighting" && q && battlePicked.value === undefined) {
    if (/^[1-9]$/.test(e.key)) { const i = Number(e.key) - 1; if (q.options && q.options[i]) { e.preventDefault(); battleAnswer(q.options[i]); } }
  } else if (battlePhase.value === "fighting" && battlePicked.value !== undefined && e.key === "Enter") {
    e.preventDefault(); battleNext();
  }
}

// ---------- 段位页 ----------
export async function openBattleHome() {
  view.value = "battle-home";
  try { battleRank.value = await api("/battle/rank?limit=20"); } catch (e) { fail(e); }
}

export async function openBattleHistory() {
  try { battleHistory.value = await api("/battle/history") || []; } catch (e) { fail(e); }
}
