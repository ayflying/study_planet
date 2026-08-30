// Package battle 段位与奖杯规则（全项目唯一实现，battle 引擎与 logic 层共用语义）。
package battle

// TierName 奖杯数 → 段位（名称 + emoji）。
// 青铜 0+ / 白银 30+ / 黄金 60+ / 铂金 100+ / 钻石 150+ / 星耀 220+ / 王者 300+。
func TierName(trophies int) (string, string) {
	switch {
	case trophies >= 300:
		return "王者", "👑"
	case trophies >= 220:
		return "星耀", "🌟"
	case trophies >= 150:
		return "钻石", "💎"
	case trophies >= 100:
		return "铂金", "🛡️"
	case trophies >= 60:
		return "黄金", "🏅"
	case trophies >= 30:
		return "白银", "🥈"
	default:
		return "青铜", "🥉"
	}
}

// TrophiesDelta 对战奖杯增减：胜 +20、平 +5、负 -10（下限 0 由调用方保证）。
func TrophiesDelta(result string) int {
	switch result {
	case "win":
		return 20
	case "draw":
		return 5
	default:
		return -10
	}
}
