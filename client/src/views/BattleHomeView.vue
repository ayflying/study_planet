<script setup>
// 段位主页：竞技场入口、段位排行榜、个人战绩
import {
  view, battleRank, battleHistory, currentStudent
} from "../state.js";
import { startBattle, openBattleHistory } from "../composables/useBattle.js";
</script>

<template>
  <section class="battle-home">
    <div class="lesson-top"><button class="close" @click="view='home'">×</button><div><p class="eyebrow orange">BATTLE ARENA</p><h2 class="page-title">对战竞技场</h2></div></div>
    <section class="card arena-hero"><div class="arena-trophy">🏆</div><div><h2>{{ battleRank?.my ? `${battleRank.my.tier_emoji} ${battleRank.my.tier}` : '🥉 青铜' }}</h2><p><b>{{ battleRank?.my?.trophies ?? 0 }}</b> 奖杯 · {{ battleRank?.my?.wins || 0 }} 胜 {{ battleRank?.my?.battles || 0 }} 战</p></div><button class="cta" @click="view='battle'; startBattle()">开始匹配对战</button></section>
    <section class="card rank-board"><h3>段位排行榜</h3><ol class="board-list"><li v-for="(e, i) in battleRank?.entries || []" :key="e.child_id" :class="{ me: e.child_id === currentStudent, top3: i < 3 }"><span class="board-rank">{{ ['🥇','🥈','🥉'][i] || e.rank }}</span><span class="board-avatar">{{ e.avatar || '🐣' }}</span><span class="board-name">{{ e.name }}<small>{{ e.tier_emoji }} {{ e.tier }}</small></span><b class="board-trophy">🏆 {{ e.trophies }}</b></li><li v-if="!(battleRank?.entries || []).length" class="board-empty">还没有人上榜，来当第一个王者！</li></ol></section>
    <section class="card history-box"><div class="history-head"><h3>我的战绩</h3><button class="history-toggle" @click="openBattleHistory">{{ battleHistory.length ? '刷新' : '加载' }}</button></div><div v-for="h in battleHistory" :key="h.id" class="history-row" :class="h.result"><span class="history-badge">{{ { win: '胜', lose: '负', draw: '平' }[h.result] }}</span><span class="board-avatar">{{ h.opponent_avatar || '🤖' }}</span><span class="history-opp">{{ h.opponent }}<small>{{ h.created_at?.slice(5, 16) }}</small></span><b class="history-score">{{ h.my_score }} : {{ h.opp_score }}</b><span class="history-trophy" :class="h.result">{{ h.trophies > 0 ? '+' : '' }}{{ h.trophies }}🏆</span></div><p v-if="!battleHistory.length" class="muted" style="text-align:center">还没有对战记录，来一场吧！</p></section>
  </section>
</template>
