// 宠物模式：打开宠物页、投喂、改名、心情表情
import {
  view, pet, petFoods, petFeeding, petMsg, showPetRename, petNewName, hasStudent
} from "../state.js";
import { api, fail } from "./useApi.js";

export async function openPet() {
  if (!hasStudent.value) return;
  view.value = "pet";
  try {
    const [p, foods] = await Promise.all([api("/pet"), api("/pet/foods")]);
    pet.value = p; petFoods.value = foods || [];
  } catch (e) { fail(e); }
}

export async function feedPet(f) {
  if (petFeeding.value || !pet.value) return;
  petFeeding.value = true; petMsg.value = "";
  try {
    const r = await api("/pet/feed", { method: "POST", body: { food: f.id } });
    pet.value = r.pet; petMsg.value = r.message || "开饭啦～";
    if (r.level_up) petMsg.value = `🎉 升到 Lv.${pet.value.level}！${petMsg.value}`;
    if (r.fed_burst) petMsg.value = `💖 好感度爆棚！${petMsg.value}`;
  } catch (e) { petMsg.value = e.message || "投喂失败"; } finally {
    petFeeding.value = false;
    setTimeout(() => { if (petMsg.value) petMsg.value = ""; }, 2600);
  }
}

export async function renamePet() {
  if (!petNewName.value.trim()) return;
  try {
    pet.value = await api("/pet/rename", { method: "POST", body: { name: petNewName.value.trim() } });
    showPetRename.value = false; petMsg.value = "改名成功～";
  } catch (e) { fail(e); }
}

export function moodEmoji(p) { return { happy: "😄", normal: "🙂", hungry: "😣", sad: "🥺" }[p?.mood] || "🙂"; }
