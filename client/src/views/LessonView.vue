<script setup>
// 练习页：题目作答（选择/填空）、判分反馈、连击特效、进度条
import { view, questions, questionIndex, combo, fx, result } from "../state.js";
import { activeQuestion, progress, selectOption, submitFill, nextQuestion, backHome } from "../composables/useLesson.js";

// 学科挑战标题映射（单一实现，禁止在模板里重复三元链）
const subjectTitles = {
  english: "英语挑战", chinese: "语文挑战", math: "数学挑战", physics: "物理挑战",
  chemistry: "化学挑战", biology: "生物挑战", history: "历史挑战", geography: "地理挑战"
};
function subjectTitle(s) { return subjectTitles[s] || "地理挑战"; }
</script>

<template>
  <section class="lesson">
    <div class="lesson-top"><button class="close" @click="view='home'">×</button><div class="progress"><i :style="{ width: `${progress}%` }" /></div><b v-if="combo >= 2">连击 x{{ combo }}</b></div>
    <article v-if="activeQuestion" class="card question" :class="{ 'review-q': activeQuestion.review }">
      <p class="eyebrow orange">{{ subjectTitle(activeQuestion.subject) }} <span v-if="activeQuestion.review" class="review-tag">错题巩固</span></p>
      <h2>{{ activeQuestion.title }}</h2>
      <p v-if="activeQuestion.passage" class="passage">{{ activeQuestion.passage }}</p>
      <em>{{ activeQuestion.subtitle }}</em>
      <div class="options">
        <button v-for="o in activeQuestion.options" :key="o" :disabled="activeQuestion.picked !== undefined" :class="{ right: activeQuestion.picked !== undefined && o === activeQuestion.rightAnswer, wrong: activeQuestion.picked === o && !activeQuestion.correct, kbdSel: activeQuestion.fillSel !== undefined && o === activeQuestion.options[activeQuestion.fillSel] }" @click="selectOption(o)">{{ o }}<i class="kbd-num">{{ activeQuestion.options.indexOf(o) + 1 }}</i></button>
      </div>
      <div v-if="activeQuestion.fillOpen && activeQuestion.picked === undefined" class="fill-box">
        <input v-model="activeQuestion.fillInput" class="fill-input" :placeholder="activeQuestion.fillHint || '输入答案，回车提交'" @keydown.enter.prevent="submitFill(activeQuestion.fillInput)" />
        <button class="fill-go" @click="submitFill(activeQuestion.fillInput)">提交 ⏎</button>
      </div>
      <div v-if="activeQuestion.picked !== undefined && activeQuestion.fillOpen" class="fill-verdict" :class="{ ok: activeQuestion.correct, no: !activeQuestion.correct }">{{ activeQuestion.correct ? '✓ 答对了' : `✗ 正确答案：${activeQuestion.rightAnswer}` }}</div>
      <p class="kbd-hint">⌨️ 数字键选择 · <b>回车键</b> 确认 / 下一题</p>
      <button v-if="activeQuestion.picked !== undefined" class="cta continue" @click="nextQuestion">{{ questionIndex + 1 < questions.length ? '继续下一题 →' : '完成关卡 ✓' }}</button>
    </article>
    <div v-if="fx" class="fx-layer">
      <template v-if="fx==='right'"><span class="fx-pop">+{{ activeQuestion?.xp || 0 }} XP</span></template>
      <template v-else-if="fx==='combo'"><span class="fx-combo">🔥 连击 x{{ combo }}！<small>+{{ activeQuestion?.xp || 0 }} XP</small></span></template>
      <template v-else-if="fx==='wrong'"><span class="fx-wrong">💪 再想想，已加入错题本</span></template>
      <template v-else-if="fx==='star'"><span class="fx-stars">★★★</span></template>
    </div>
    <div v-if="result" class="modal"><section class="result card"><div class="result-star">★</div><p class="eyebrow orange">MISSION COMPLETE</p><h2>闯关完成！</h2><p>答对 {{ result.correct }}/{{ result.total }} 题 · 最高连击 x{{ result.max_combo }}</p><div class="result-stars">{{ '★'.repeat(result.stars || 0) }}<i>{{ '☆'.repeat(3 - (result.stars || 0)) }}</i></div><div v-if="result.xp_gained" class="result-xp">✨ 经验值 +{{ result.xp_gained }}</div><button class="cta" @click="backHome">回到学习地图</button></section></div>
  </section>
</template>
