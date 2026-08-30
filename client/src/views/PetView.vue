<script setup>
// 宠物页：宠物状态卡（饱食度/好感度/经验条）+ 投喂零食
import {
  view, pet, petFoods, petFeeding, petMsg, showPetRename, petNewName
} from "../state.js";
import { feedPet, renamePet, moodEmoji } from "../composables/usePet.js";
</script>

<template>
  <section class="pet-page">
    <div class="lesson-top"><button class="close" @click="view='home'">×</button><div><p class="eyebrow orange">MY PET</p><h2 class="page-title">我的宠物</h2></div></div>
    <section v-if="pet" class="card pet-card" :class="`mood-${pet.mood}`"><div class="pet-stage"><div class="pet-bubble">{{ petMsg || (pet.mood === 'hungry' ? '咕噜咕噜…肚子饿了~' : pet.mood === 'sad' ? '好久没陪我玩了…' : '和我一起学习吧！') }}</div><div class="pet-avatar" :class="{ hop: petMsg }">{{ pet.emoji || '🐣' }}</div><button class="pet-rename" @click="showPetRename=true; petNewName=pet.name">✏️</button></div><div class="pet-name-row"><b>{{ pet.name }}</b><span>{{ pet.species_name }}</span><em>{{ moodEmoji(pet) }} {{ pet.mood_text }}</em></div><div class="pet-bars"><div class="pet-bar"><label>饱食度</label><div class="bar"><i :style="{ width: `${pet.hunger}%` }" class="bar-hunger"/></div><b>{{ pet.hunger }}%</b></div><div class="pet-bar"><label>好感度</label><div class="bar"><i :style="{ width: `${pet.affection}%` }" class="bar-aff"/></div><b>{{ pet.affection }}%</b></div><div class="pet-bar"><label>Lv.{{ pet.level }}</label><div class="bar"><i :style="{ width: `${Math.round(pet.exp / pet.exp_max * 100)}%` }" class="bar-exp"/></div><b>{{ pet.exp }}/{{ pet.exp_max }}</b></div></div></section>
    <section v-if="pet" class="card food-box"><h3>🍖 投喂零食</h3><p class="muted">不同食物增加的饱食度/好感度/经验不同，饱食度会随时间下降，记得常回来喂它！</p><div class="food-grid"><button v-for="f in petFoods" :key="f.id" class="food" :disabled="petFeeding" @click="feedPet(f)"><span class="food-emoji">{{ f.emoji }}</span><b>{{ f.name }}</b><small>🍚 +{{ f.hunger }} 💗 +{{ f.affection }} ✨ +{{ f.exp }}</small></button></div></section>
  </section>
</template>
