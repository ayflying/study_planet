<script setup>
import { computed, onMounted, ref } from "vue";

const API = "/api";
const students = ref([]);
const currentStudent = ref(Number(localStorage.getItem("sp_stu") || 1));
const points = ref(0);
const sessions = ref([]);
const words = ref([]);
const math = ref([]);
const reading = ref(null);
const view = ref("home");
const lesson = ref(null);
const questions = ref([]);
const questionIndex = ref(0);
const combo = ref(0);
const sessionId = ref(null);
const result = ref(null);
const message = ref("");
const loading = ref(true);

const units = [
  { kind: "words", title: "单词关 · 第1关", sub: "认识 10 个单词", className: "word" },
  { kind: "reading", title: "阅读关 · 小蚂蚁搬骨头", sub: "读短文答问题", className: "read" },
  { kind: "math", title: "数学关 · 第1题组", sub: "计算与应用", className: "math" }
];

const activeQuestion = computed(() => questions.value[questionIndex.value]);
const progress = computed(() => questions.value.length ? Math.round(questionIndex.value / questions.value.length * 100) : 0);
const petMood = computed(() => combo.value >= 3 ? "wow" : "idle");

function api(path, options = {}) {
  const separator = path.includes("?") ? "&" : "?";
  return fetch(`${API}${path}${separator}student_id=${currentStudent.value}`, {
    method: options.method || "GET",
    headers: { "Content-Type": "application/json" },
    body: options.body ? JSON.stringify(options.body) : undefined
  }).then(async response => {
    const data = await response.json().catch(() => null);
    if (!response.ok) throw new Error(data?.error || `HTTP ${response.status}`);
    return data;
  });
}

function shuffle(items) {
  return [...items].sort(() => Math.random() - 0.5);
}

function starCount(unit) {
  return sessions.value.filter(item => item.subject === unit.kind && item.level === 1)
    .reduce((best, item) => Math.max(best, item.stars || 0), 0);
}

function totalStars() {
  return sessions.value.reduce((sum, item) => sum + (item.stars || 0), 0);
}

function maxCombo() {
  return sessions.value.reduce((best, item) => Math.max(best, item.max_combo || 0), 0);
}

function petSvg(mood) {
  const mouth = mood === "wow"
    ? '<ellipse cx="60" cy="72" rx="8" ry="10" fill="#7A3B12"/><ellipse cx="60" cy="76" rx="4.5" ry="4.5" fill="#FF8FAF"/>'
    : '<path d="M51 69 q9 8 18 0" stroke="#7A3B12" stroke-width="3.5" fill="none" stroke-linecap="round"/>';
  return `<svg viewBox="0 0 120 120" width="112" height="112" aria-label="学习伙伴">
    <ellipse cx="60" cy="72" rx="36" ry="34" fill="#FFCE54"/><ellipse cx="60" cy="80" rx="24" ry="20" fill="#FFE291"/>
    <polygon points="30,48 26,24 48,40" fill="#FFCE54" stroke="#E8A400" stroke-linejoin="round" stroke-width="3"/><polygon points="90,48 94,24 72,40" fill="#FFCE54" stroke="#E8A400" stroke-linejoin="round" stroke-width="3"/>
    <circle cx="47" cy="57" r="5.5" fill="#3A2E25"/><circle cx="73" cy="57" r="5.5" fill="#3A2E25"/><circle cx="49" cy="55" r="1.8" fill="#fff"/><circle cx="75" cy="55" r="1.8" fill="#fff"/>
    <ellipse cx="36" cy="68" rx="6" ry="4" fill="#FFC4B8"/><ellipse cx="84" cy="68" rx="6" ry="4" fill="#FFC4B8"/>${mouth}
    <ellipse cx="24" cy="84" rx="8" ry="6" fill="#FFCE54"/><ellipse cx="96" cy="84" rx="8" ry="6" fill="#FFCE54"/>
  </svg>`;
}

function optionPool(answer, pool) {
  const options = [answer];
  shuffle(pool).forEach(item => { if (item && !options.includes(item) && options.length < 4) options.push(item); });
  while (options.length < 4) options.push(`选项${options.length}`);
  return shuffle(options.map(String));
}

async function load() {
  loading.value = true;
  try {
    const [studentData, pointData, wordData, mathData, sessionData] = await Promise.all([
      api("/students"), api("/points"), api("/words"), api("/math"), api("/sessions")
    ]);
    students.value = studentData || [];
    points.value = pointData?.total || 0;
    words.value = wordData || [];
    math.value = mathData || [];
    sessions.value = sessionData || [];
  } catch (error) {
    message.value = `连不上服务器：${error.message}`;
  } finally {
    loading.value = false;
  }
}

async function switchStudent(id) {
  currentStudent.value = id;
  localStorage.setItem("sp_stu", String(id));
  await load();
}

async function prepareReading() {
  if (!reading.value) reading.value = await api("/readings/1");
}

async function startUnit(unit) {
  try {
    if (unit.kind === "reading") await prepareReading();
    let list = [];
    if (unit.kind === "words") {
      list = shuffle(words.value).slice(0, 5).map(word => ({
        id: word.id, type: "word", title: word.word, subtitle: word.phonetic, answer: word.meaning,
        options: optionPool(word.meaning, words.value.map(item => item.meaning))
      }));
    }
    if (unit.kind === "math") {
      list = shuffle(math.value).slice(0, 5).map(item => ({
        id: item.id, type: "math", title: item.question, explain: item.explanation, answer: String(item.answer),
        options: shuffle(JSON.parse(item.options || "[]").map(String))
      }));
    }
    if (unit.kind === "reading") {
      list = shuffle(reading.value?.questions || []).slice(0, 5).map(item => {
        const raw = String(item.answer || "");
        const options = [item.option_a, item.option_b, item.option_c, item.option_d].filter(Boolean).map(stripPrefix);
        const answer = /^[A-D]$/.test(raw) ? options["ABCD".indexOf(raw)] : stripPrefix(raw);
        return { id: item.id, type: "reading", title: item.question, answer, options: shuffle(options) };
      });
    }
    if (!list.length) throw new Error("这一关还没有题目，让家长先添加内容吧");
    const created = await api("/sessions", { method: "POST", body: { subject: unit.kind, level: 1, total: list.length } });
    lesson.value = unit;
    questions.value = list;
    questionIndex.value = 0;
    combo.value = 0;
    sessionId.value = created.id;
    view.value = "lesson";
    message.value = "准备好闯关了吗？";
  } catch (error) {
    message.value = error.message;
  }
}

function stripPrefix(value) { return String(value || "").replace(/^[A-D][.、\s]\s*/, ""); }

async function selectOption(option) {
  const question = activeQuestion.value;
  if (!question || question.picked) return;
  question.picked = option;
  question.correct = option === question.answer;
  try {
    let response;
    if (question.type === "word") response = await api(`/words/${question.id}/progress`, { method: "POST", body: { known: question.correct, session_id: sessionId.value } });
    else if (question.type === "math") response = await api(`/math/${question.id}/answer`, { method: "POST", body: { answer: option, session_id: sessionId.value } });
    else response = await api("/readings/1/answer", { method: "POST", body: { question_id: question.id, answer: option, session_id: sessionId.value } });
    question.correct = response.correct ?? question.correct;
    combo.value = question.correct ? (response.combo || combo.value + 1) : 0;
    message.value = question.correct ? (combo.value >= 3 ? `太棒了！连击 x${combo.value}` : "答对了哦！") : (question.explain ? `解析：${question.explain}` : "没关系，再想想～");
  } catch {
    combo.value = 0;
    message.value = "提交失败，请检查网络后重试。";
  }
}

async function nextQuestion() {
  if (questionIndex.value + 1 < questions.value.length) { questionIndex.value += 1; return; }
  try {
    result.value = await api(`/sessions/${sessionId.value}/finish`, { method: "POST" });
    const pointData = await api("/points");
    points.value = pointData?.total || points.value;
    sessions.value = await api("/sessions");
  } catch (error) { message.value = `结算失败：${error.message}`; }
}

function backHome() {
  result.value = null;
  view.value = "home";
  message.value = "我在这儿陪你学习～";
}

onMounted(load);
</script>

<template>
  <main class="shell">
    <header class="topbar">
      <div class="brand"><div class="logo"><img src="/assets/logo.png" alt="学霸星球 logo" /></div><div><h1>学霸星球</h1><p>今天也要闯关哦</p></div></div>
      <div class="xp-ring"><strong>{{ points }}</strong><span>积分</span></div>
    </header>

    <section v-if="loading" class="card loading">正在装载学习星球…</section>
    <section v-else-if="view === 'home'">
      <div class="students"><button v-for="student in students" :key="student.id" class="student" :class="{ active: student.id === currentStudent }" @click="switchStudent(student.id)"><span>{{ student.avatar || '🐣' }}</span>{{ student.name }}</button></div>
      <div class="path card"><h2>🌎 学习星球航线</h2><button v-for="unit in units" :key="unit.kind" class="unit" :class="unit.className" @click="startUnit(unit)"><span class="stars">{{ '★'.repeat(starCount(unit)) }}<i>{{ '☆'.repeat(3 - starCount(unit)) }}</i></span><strong>{{ unit.title }}</strong><small>{{ unit.sub }} · 点击开始</small></button></div>
      <div class="pet-zone"><div class="bubble">{{ message || '准备好闯关了吗？' }}</div><div class="pet" :class="petMood" v-html="petSvg(petMood)" /></div>
      <div class="stats"><div><b>{{ points }}</b><span>我的积分</span></div><div><b>{{ maxCombo() }}</b><span>最高连击</span></div><div><b>{{ totalStars() }}</b><span>总星星</span></div></div>
    </section>

    <section v-else class="lesson">
      <div class="lesson-top"><button class="close" @click="view = 'home'">×</button><div class="progress"><i :style="{ width: `${progress}%` }" /></div><b v-if="combo >= 2">🔥x{{ combo }}</b></div>
      <article v-if="activeQuestion" class="card question"><p>{{ { word: '听一听 · 选释义', math: '算一算', reading: '读一读 · 选一选' }[activeQuestion.type] }}</p><h2>{{ activeQuestion.title }}</h2><em v-if="activeQuestion.subtitle">{{ activeQuestion.subtitle }}</em><div class="options"><button v-for="option in activeQuestion.options" :key="option" :disabled="activeQuestion.picked" :class="{ right: activeQuestion.picked && option === activeQuestion.answer, wrong: activeQuestion.picked === option && !activeQuestion.correct }" @click="selectOption(option)">{{ option }}</button></div><button v-if="activeQuestion.picked" class="continue" @click="nextQuestion">{{ questionIndex + 1 < questions.length ? '继续 →' : '完成关卡 ✓' }}</button></article>
      <div class="pet-zone compact"><div class="bubble">{{ message }}</div><div class="pet" v-html="petSvg(activeQuestion?.correct ? 'wow' : 'idle')" /></div>
    </section>

    <div v-if="result" class="modal"><section class="result card"><div class="big-stars">{{ '★'.repeat(result.stars || 0) }}<i>{{ '☆'.repeat(3 - (result.stars || 0)) }}</i></div><h2>{{ ['再接再厉！', '不错哦！', '太棒了！', '完美通关！'][result.stars || 0] }}</h2><p>答对 {{ result.correct }}/{{ result.total }} 题 · 最高连击 x{{ result.max_combo }}</p><div class="stats"><div><b>{{ result.stars }}/3</b><span>星星</span></div><div><b>x{{ result.max_combo }}</b><span>连击</span></div><div><b>+{{ result.bonus }}</b><span>奖分</span></div></div><button class="continue" @click="backHome">回星球 🏠</button></section></div>
  </main>
</template>
