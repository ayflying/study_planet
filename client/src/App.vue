<script setup>
// App.vue：应用壳。只负责顶栏、视图路由（view 状态分发到各视图组件）、底部菜单与全局初始化。
// 业务逻辑分别在 composables/（useApi/useAuth/useData/useLesson/useBattle/usePet/useAdmin），
// 视图在 views/，共享状态在 state.js，常量在 constants.js。
import { onMounted } from "vue";
import {
  view, loading, students, inAdmin, adminTab, canAdmin, parentToken,
  hasStudent, points
} from "./state.js";
import { api } from "./composables/useApi.js";
import { enterAdmin, backToStudent, loadAuth, markAdmin } from "./composables/useAuth.js";
import { load, loadWrong } from "./composables/useData.js";
import { lessonKeyHandler, lessonNumHandler } from "./composables/useLesson.js";
import { battleKeyHandler } from "./composables/useBattle.js";
import { loadAdmin } from "./composables/useAdmin.js";

import HomeView from "./views/HomeView.vue";
import LessonView from "./views/LessonView.vue";
import BattleView from "./views/BattleView.vue";
import BattleHomeView from "./views/BattleHomeView.vue";
import PetView from "./views/PetView.vue";
import AdminView from "./views/AdminView.vue";
import GlobalModals from "./views/GlobalModals.vue";

onMounted(async () => {
  window.addEventListener("keydown", e => { lessonKeyHandler(e); battleKeyHandler(e); lessonNumHandler(e); });
  await loadAuth(); await load();
  if (parentToken.value) {
    markAdmin();
    await loadAdmin().catch(() => {});
    view.value = "admin";
  }
  await loadWrong();
});
</script>

<template>
  <main class="shell">
    <header class="topbar"><div class="brand"><div class="logo"><img src="/logo.png" alt="学霸星球" /></div><div><p class="eyebrow">STUDY PLANET · 学习成长空间</p><h1>学霸星球</h1><p>{{ inAdmin ? '家长管理中心' : '专为孩子设计的学习旅程' }}</p></div></div><div class="top-actions"><div v-if="hasStudent" class="xp-ring"><strong>{{ points }}</strong><span>积分</span></div><button v-if="!inAdmin" class="parent-entry" @click="enterAdmin()">家长中心</button><button v-else class="parent-entry" @click="backToStudent">返回学习</button></div></header>
    <section v-if="loading" class="card loading">正在装载学习星球…</section>
    <section v-else-if="!students.length" class="welcome card"><div class="welcome-art"><div class="planet">学</div><span class="orbit orbit-a"></span><span class="orbit orbit-b"></span></div><div class="welcome-copy"><p class="eyebrow orange">家长先登录，开启学习旅程</p><h2>先管理家庭成员<br />再让孩子开始学习</h2><p>家长登录后，可以创建孩子账号、安排每日任务、设置奖励，并查看学习成长记录。</p><button class="cta" @click="enterAdmin()">家长认证进入 <span>→</span></button></div><div class="feature-strip"><span>✦ 安全的家庭空间</span><span>✦ 趣味化闯关学习</span><span>✦ 可追踪成长记录</span></div></section>
    <HomeView v-else-if="view === 'home'" />
    <LessonView v-else-if="view === 'lesson'" />
    <BattleView v-else-if="view === 'battle'" />
    <BattleHomeView v-else-if="view === 'battle-home'" />
    <PetView v-else-if="view === 'pet'" />
    <AdminView v-else />
    <nav v-if="inAdmin && view !== 'lesson'" class="bottom-menu" aria-label="主菜单"><button :class="{ active: view === 'home' }" @click="view = 'home'"><span class="menu-icon">⌂</span><span>学习</span></button><button :class="{ active: view === 'admin' && adminTab === 'tasks' }" @click="enterAdmin('tasks')"><span class="menu-icon">✓</span><span>任务</span></button><button :class="{ active: view === 'admin' && adminTab === 'rewards' }" @click="enterAdmin('rewards')"><span class="menu-icon">★</span><span>奖励</span></button><button :class="{ active: view === 'admin' && ['overview', 'students', 'logs', 'settings'].includes(adminTab) }" @click="enterAdmin('overview')"><span class="menu-icon">♙</span><span>家长</span></button></nav>
    <GlobalModals />
  </main>
</template>
