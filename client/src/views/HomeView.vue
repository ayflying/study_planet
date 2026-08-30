<script setup>
// 首页：学习地图 + 快捷入口 + 错题巩固 + 成长统计
import {
  hasStudent, students, currentStudent, canAdmin, activeStudent, points, wrongCount
} from "../state.js";
import { enterAdmin, switchStudent } from "../composables/useAuth.js";
import { units, wrongSubject, starCount, totalStars, maxCombo } from "../composables/useData.js";
import { startUnit } from "../composables/useLesson.js";
import { openBattleHome } from "../composables/useBattle.js";
import { openPet } from "../composables/usePet.js";
import { openMyTasks, openMyRewards } from "../composables/useAdmin.js";
import { showLeaderboard } from "../composables/useData.js";
</script>

<template>
  <section class="home">
    <div v-if="!hasStudent" class="empty-state card"><div class="empty-icon">👨‍👩‍👧‍👦</div><p class="eyebrow orange">第一步 · 创建学习账号</p><h2>还没有学生账号</h2><p>请先在家长中心添加孩子，完成后这里就会出现专属学习工作台。</p><button class="cta" @click="enterAdmin('students')">去添加学生 <span>→</span></button></div>
    <template v-else><section class="dashboard-hero"><div><p class="eyebrow orange">{{ activeStudent?.name }}的学习星球</p><h2>今天想去哪里探索？</h2><p>每一次练习，都会变成成长路上的一颗星。</p></div><div class="hero-badge">{{ activeStudent?.avatar || '🐣' }}<small>Lv.{{ activeStudent?.level || 1 }}</small></div></section><div class="student-switcher"><span class="section-label">当前学习者</span><button v-for="s in students" :key="s.id" class="student" :class="{ active: s.id === currentStudent }" @click="switchStudent(s.id)"><span>{{ s.avatar || '🐣' }}</span>{{ s.name }}<b v-if="s.id === currentStudent">✓</b></button><button v-if="canAdmin" class="add-child" @click="enterAdmin('students')">＋ 添加孩子</button></div><section class="path card"><div class="section-heading"><div><p class="eyebrow">LEARNING MAP</p><h2>学习探索地图</h2></div><span class="map-progress">{{ totalStars() }} 颗星星</span></div><div class="units"><button v-for="u in units" :key="u.kind" class="unit" :class="u.className" @click="startUnit(u)"><span class="unit-icon">{{ u.icon }}</span><span class="stars">{{ '★'.repeat(starCount(u)) }}<i>{{ '☆'.repeat(3 - starCount(u)) }}</i></span><strong>{{ u.title }}</strong><small>{{ u.sub }}</small><em>开始探索 <b>→</b></em></button></div></section><div class="quick-grid"><button @click="openBattleHome"><span class="quick-icon rank">⚔️</span><span><b>真人对战</b><small>5题抢答 · 冲击王者段位</small></span><i>→</i></button><button @click="openMyTasks"><span class="quick-icon">✓</span><span><b>我的任务</b><small>完成任务赢积分</small></span><i>→</i></button><button @click="openMyRewards"><span class="quick-icon gift">◆</span><span><b>奖励商店</b><small>用积分兑换奖励</small></span><i>→</i></button><button @click="openPet"><span class="quick-icon pet-icon">🐣</span><span><b>我的宠物</b><small>投喂零食养好感度</small></span><i>→</i></button><button @click="showLeaderboard"><span class="quick-icon rank2">♛</span><span><b>本周排行榜</b><small>比一比谁的经验多</small></span><i>→</i></button></div><section v-if="wrongCount" class="wrong-card card" @click="startUnit({ kind: wrongSubject(), title: '错题巩固' })"><div class="wrong-icon">📕</div><div class="wrong-copy"><b>错题本有 {{ wrongCount }} 道题待巩固</b><small>练习时会自动穿插错题，答对即可消除</small></div><i>→</i></section><section class="growth card"><div class="pet-zone"><div class="bubble">准备好闯关了吗？</div><div class="pet">🌟</div></div><div class="stats"><div><b>{{ points }}</b><span>我的积分</span></div><div><b>{{ maxCombo() }}</b><span>最高连击</span></div><div><b>{{ totalStars() }}</b><span>收集星星</span></div></div></section></template>
  </section>
</template>
