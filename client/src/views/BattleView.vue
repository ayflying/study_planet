<script setup>
// 真人对战页：匹配中 / 对战中（HUD+计时器+抢答）/ 结算画面
import {
  view, battlePhase, battleMe, battleOpp, battleQs, battleIdx, battleRemain,
  battlePicked, battleResult, battleFillInput, activeStudent
} from "../state.js";
import { startBattle, battleAnswer, exitBattle, openBattleHome } from "../composables/useBattle.js";
import { gradeLabel } from "../constants.js";
</script>

<template>
  <section class="battle">
    <div v-if="battlePhase === 'matching'" class="card battle-match"><div class="vs-orbit"><span class="vs-me">{{ activeStudent?.avatar || '🐣' }}</span><span class="vs-dot"></span><span class="vs-q">?</span></div><h2>正在匹配对手…</h2><p class="muted">真人玩家 · 数学 · {{ gradeLabel(activeStudent?.grade) }}<br />3 秒内没找到会安排机器人陪练</p><div class="match-pulse"><i></i><i></i><i></i></div><button class="quit-link" @click="exitBattle">取消匹配</button></div>
    <!-- 对战中 -->
    <template v-else-if="battlePhase === 'fighting'">
      <div class="battle-hud card"><div class="hud-side me"><span class="hud-avatar">{{ activeStudent?.avatar || '🐣' }}</span><b>{{ battleMe.total }}</b><small>{{ activeStudent?.name }}</small></div><div class="hud-mid"><span class="hud-qnum">{{ battleIdx + 1 }}/{{ battleQs.length }}</span><div class="hud-timer" :class="{ danger: battleRemain <= 3 }"><svg viewBox="0 0 80 80"><circle cx="40" cy="40" r="34" class="timer-bg"/><circle cx="40" cy="40" r="34" class="timer-fg" :style="{ strokeDashoffset: 214 - 214 * battleRemain / 10 }"/></svg><b>{{ battleRemain }}</b></div></div><div class="hud-side opp"><span class="hud-avatar">{{ battleOpp.avatar || '🤖' }}</span><b>{{ battleOpp.total }}</b><small>{{ battleOpp.name }}{{ battleOpp.isBot ? ' 🤖' : '' }}</small></div></div>
      <article v-if="battleQs[battleIdx]" class="card question battle-q">
        <p class="eyebrow orange">{{ battleQs[battleIdx].topic || '抢答' }} · {{ { choice: '选择题', fill: '填空题' }[battleQs[battleIdx].qtype] || '抢答' }}</p>
        <h2>{{ battleQs[battleIdx].question }}</h2>
        <div v-if="battleQs[battleIdx].options?.length" class="options battle-options" :class="{ done: battlePicked !== undefined }"><button v-for="o in battleQs[battleIdx].options" :key="o" :disabled="battlePicked !== undefined" @click="battleAnswer(o)">{{ o }}<i class="kbd-num">{{ battleQs[battleIdx].options.indexOf(o) + 1 }}</i></button></div>
        <div v-else class="fill-box"><input v-model="battleFillInput" class="fill-input" placeholder="输入答案，回车抢答" @keydown.enter.prevent="battleAnswer(battleFillInput); battleFillInput=''" /><button class="fill-go" @click="battleAnswer(battleFillInput); battleFillInput=''">抢答 ⏎</button></div>
        <div v-if="battlePicked !== undefined && !battleCorrect" class="battle-verdict no">✗ 这题没拿到分<small>当前 {{ battleMe.total }} : {{ battleOpp.total }}</small></div>
      </article>
    </template>
    <!-- 结算画面 -->
    <div v-else-if="battlePhase === 'finished' && battleResult" class="modal" style="position:static;display:grid;place-items:center">
      <section class="card battle-result" :class="battleResult.result">
        <template v-if="battleResult.result === 'win'"><div class="result-crown">👑</div><h2>胜利！</h2></template>
        <template v-else-if="battleResult.result === 'lose'"><div class="result-crown">💪</div><h2>惜败</h2></template>
        <template v-else><div class="result-crown">🤝</div><h2>平局</h2></template>
        <div class="battle-score"><span>{{ battleResult.my_score }}</span><i>:</i><span>{{ battleResult.opp_score }}</span></div>
        <div class="battle-trophies"><b>{{ battleResult.trophies }}</b><span>奖杯</span><em>{{ battleResult.tier_emoji }} {{ battleResult.tier }}</em></div>
        <ul class="battle-rewards"><li v-for="(r, i) in battleResult.rewards || []" :key="i">{{ r }}</li></ul>
        <div v-if="battleResult.snack_name" class="battle-snack">🎁 掉落零食：{{ battleResult.snack_name }} ×1</div>
        <div class="battle-exp">✨ 经验 +{{ battleResult.exp }}</div>
        <button class="cta" @click="exitBattle(); openBattleHome()">查看段位榜</button>
        <button class="quit-link" @click="exitBattle(); view='home'">返回主页</button>
      </section>
    </div>
  </section>
</template>
