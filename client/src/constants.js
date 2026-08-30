// 全局共享常量：只定义一处，禁止在组件里硬编码魔法值
export const API_BASE = "/api";
export const WS_URL = `${location.protocol === "https:" ? "wss" : "ws"}://${location.host}/ws/battle`;

export const gradeOptions = [
  { value: 1, label: "一年级" },
  { value: 2, label: "二年级" },
  { value: 3, label: "三年级" },
  { value: 4, label: "四年级" },
  { value: 5, label: "五年级" },
  { value: 6, label: "六年级" },
  { value: 7, label: "初一" },
  { value: 8, label: "初二" },
  { value: 9, label: "初三" }
];

export const avatarOptions = ["🐣", "🚀", "🌟", "🦊", "🐼", "🦄", "🐬", "🐨"];

// 学习地图单元配色（循环使用）
export const unitColors = ["word", "read", "math", "sci", "sci2", "bio", "his", "geo"];

// 学科兜底列表（subjects 接口异常时使用）
export const fallbackUnits = [
  { code: "english", name: "英语", icon: "Aa", count: 0 },
  { code: "chinese", name: "语文", icon: "文", count: 0 },
  { code: "math", name: "数学", icon: "∑", count: 0 }
];

// 错题巩固时学科优先级
export const wrongSubjectOrder = ["math", "english", "chinese", "physics", "chemistry", "biology", "history", "geography"];

export function gradeLabel(g) {
  const opt = gradeOptions.find(x => x.value === (g || 1));
  return opt ? opt.label : `${g || 1}年级`;
}
