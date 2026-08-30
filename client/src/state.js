// 全局共享状态：所有跨模块 ref 集中在此，避免 composable 之间循环依赖
import { ref, computed } from "vue";

export const parentToken = ref(localStorage.getItem("sp_parent_jwt") || "");
export const authMode = ref("casdoor");
export const loading = ref(true);
// 全局请求计数：>0 时显示阻塞式 Loading 弹框，防止网慢时连点/乱点（页面首次加载除外）
export const busy = ref(0);
export const error = ref("");
export const notice = ref("");
export const view = ref("home");
export const adminTab = ref("overview");

// 权限分离：学生界面默认可进（无需家长登录），家长中心必须重新认证（Casdoor 或家长 PIN）。
// canAdmin 只在家长完成认证后为真；一旦返回学生界面立即失效，再次进入家长中心需重新认证。
// sessionStorage 保证：刷新不丢（Casdoor 回跳需要），关闭标签页后自动回到学生模式。
export const canAdmin = ref(sessionStorage.getItem("sp_admin") === "1");
export const inAdmin = computed(() => canAdmin.value && view.value === "admin");
export const pendingTab = ref("overview");

// ---------- 学生 ----------
export const students = ref([]);
export const currentStudent = ref(Number(localStorage.getItem("sp_stu") || 0));
export const points = ref(0);
export const sessions = ref([]);
export const subjects = ref([]);
export const hasStudent = computed(() => students.value.length > 0 && !!currentStudent.value);
// 学生模式：有学生列表即可进入学习工作台（学生列表公开接口拉取，不含敏感数据）
export const studentMode = computed(() => students.value.length > 0 && !canAdmin.value);
export const activeStudent = computed(() => students.value.find(s => s.id === currentStudent.value));

// ---------- 家长中心数据 ----------
export const tasks = ref([]), rewards = ref([]), logs = ref([]);
export const newTask = ref({ title: "", type: "学习", due_date: "", points: 5 });
export const newReward = ref({ name: "", cost_points: 20 });
export const newStudent = ref({ name: "", username: "", avatar: "🐣", grade: 1 });
export const showAddStudent = ref(false);

// ---------- 练习（lesson） ----------
export const lesson = ref(null), questions = ref([]), questionIndex = ref(0), combo = ref(0), sessionId = ref(null), result = ref(null);
export const fx = ref(""), leaderboard = ref(null), wrongCount = ref(0), showBoard = ref(false);
export const showMyTasks = ref(false), myTasks = ref([]);
export const showMyRewards = ref(false), myRewards = ref([]);

// ---------- 宠物模式 ----------
export const pet = ref(null), petFoods = ref([]), petFeeding = ref(false), petMsg = ref("");
export const showPetRename = ref(false), petNewName = ref("");

// ---------- 真人对战 ----------
export const battlePhase = ref("idle"); // idle | matching | fighting | finished
export const battleMe = ref({ total: 0 }), battleOpp = ref({ total: 0, name: "", avatar: "", isBot: false });
export const battleQs = ref([]), battleIdx = ref(0), battleRemain = ref(10);
export const battlePicked = ref(undefined), battleCorrect = ref(false), battleGain = ref(0);
export const battleFillInput = ref("");
export const battleResult = ref(null);
export const battleRank = ref(null), battleHistory = ref([]);
