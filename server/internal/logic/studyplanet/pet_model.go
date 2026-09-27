// Package studyplanet 宠物核心模型：种类/食物表、等级公式、心情推导、读写与惰性衰减。
package studyplanet

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	v1 "studyplanet/api/studyplanet/v1"
)

// petSpecies 宠物种类表。
type petSpecies struct {
	Code  string
	Name  string
	Emoji string
}

var petSpeciesList = []petSpecies{
	{"cat", "小猫咪", "🐱"},
	{"dog", "小狗狗", "🐶"},
	{"rabbit", "小兔子", "🐰"},
	{"panda", "小熊猫", "🐼"},
	{"fox", "小狐狸", "🦊"},
	{"unicorn", "小独角兽", "🦄"},
}

// petFoods 投喂食物表（id → 恢复值）。
type petFood struct {
	ID        string
	Name      string
	Emoji     string
	Hunger    int
	Affection int
	Exp       int
	Desc      string
}

var petFoodList = []petFood{
	{"apple", "苹果", "🍎", 10, 3, 1, "脆脆甜甜，日常小点心"},
	{"fish", "小鱼干", "🐟", 18, 5, 2, "猫猫的最爱！"},
	{"cake", "小蛋糕", "🍰", 25, 8, 3, "节日限定，好感度暴涨"},
	{"hotpot", "小火锅", "🍲", 40, 10, 5, "大餐一顿，饱食度大回复"},
	{"milk", "牛奶", "🥛", 12, 4, 1, "睡前一杯，长得高高的"},
	{"star", "星星糖", "⭐", 8, 12, 2, "不含饱但超级喜欢！好感度神器"},
}

// ---------- 零食掉落权重 ----------
// 每把最多获得1个零食，最少0个。
// 用户给定比例：苹果40%、小鱼干/牛奶/星星糖各30%、小蛋糕20%、小火锅10%。
// 比例之和160，故以权重实现：不掉落权重40，总权重200。
// 随机1..200，命中区间即得对应零食——保留用户给定的相对稀有度。
type snackDrop struct {
	Food   string
	Weight int
}

var snackDropTable = []snackDrop{
	{"", 40}, // 不掉落
	{"apple", 40},
	{"fish", 30},
	{"milk", 30},
	{"star", 30},
	{"cake", 20},
	{"hotpot", 10},
}
var snackDropTotal = 200

// rollSnack 掷一次零食掉落：命中返回 foodID，否则返回 ""。
func rollSnack() string {
	r := grand.N(1, snackDropTotal)
	for _, d := range snackDropTable {
		if r <= d.Weight {
			return d.Food
		}
		r -= d.Weight
	}
	return ""
}

// snackLabel foodID → 显示文本（emoji + 名称），用于"获得零食"提示。
func snackLabel(food string) string {
	for _, f := range petFoodList {
		if f.ID == food {
			return f.Emoji + " " + f.Name
		}
	}
	return ""
}

// inventoryColumn foodID → pets 表列名。
func inventoryColumn(food string) string {
	switch food {
	case "apple":
		return "food_apple"
	case "fish":
		return "food_fish"
	case "milk":
		return "food_milk"
	case "star":
		return "food_star"
	case "cake":
		return "food_cake"
	case "hotpot":
		return "food_hotpot"
	}
	return ""
}

// ---------- 宠物基础 ----------

// petExpNeed 升到 level+1 所需累计经验：20 + (level-1)*15。
func petExpNeed(level int) int { return 20 + (level-1)*15 }

// petHungerHours 每满24小时饱食度衰减当前值的20%（指数衰减，越饿掉得越慢）。
const petHungerHours = 24

// petHungerPct 每24小时扣除的当前饱食度比例（20%）。
const petHungerPct = 0.2

// petAffectionHours 每满8小时好感度衰减1（缓慢衰减，长期不互动才会触发惩罚）。
const petAffectionHours = 8

// petMood 由饱食度+好感度推导心情。
func petMood(hunger, affection int) (mood, text string) {
	score := float64(hunger)*0.4 + float64(affection)*0.6
	switch {
	case score >= 70:
		return "happy", "开心得转圈圈～"
	case score >= 45:
		return "normal", "安静地陪着你学习"
	case score >= 20:
		return "sad", "有点想你陪陪它…"
	default:
		return "hungry", "饿扁了，快投喂我！"
	}
}

// petOf 取宠物（不存在则自动创建）。
func (s *sStudyPlanet) petOf(ctx context.Context, childID int) (gdb.Record, error) {
	row, err := g.DB().Model("pets").Ctx(ctx).Where("child_id", childID).One()
	if err != nil {
		return nil, err
	}
	if !row.IsEmpty() {
		return row, nil
	}
	// 首次领取：随机一种宠物
	sp := petSpeciesList[grand.N(0, len(petSpeciesList)-1)]
	if _, err := g.DB().Model("pets").Ctx(ctx).Data(g.Map{
		"child_id":          childID,
		"name":              sp.Name,
		"species":           sp.Code,
		"hunger":            80,
		"affection":         20,
		"last_decay_at":     gtime.Now(),
		"last_hunger_at":    gtime.Now(),
		"last_affection_at": gtime.Now(),
		// 食物库存默认为0（由迁移 DEFAULT 0 保证）
	}).Insert(); err != nil {
		return nil, err
	}
	return g.DB().Model("pets").Ctx(ctx).Where("child_id", childID).One()
}

// petTick 惰性结算：按真实流逝时间分别衰减饱食度（每 24h 扣当前值 20%）与
// 好感度（每 8h -1）。两者使用独立时间戳（last_hunger_at / last_affection_at），
// 互不重置——否则好感度频繁衰减会不断刷新共享时间戳，导致饱食度 24h 衰减对
// 频繁打开宠物页的孩子永不触发。缺失时间戳时回退到 last_decay_at / 当前时间。
func (s *sStudyPlanet) petTick(ctx context.Context, pet gdb.Record) {
	childID := pet["child_id"].Int()
	hunger := pet["hunger"].Int()
	aff := pet["affection"].Int()

	now := gtime.Now()
	lastH := pet["last_hunger_at"].Time()
	if lastH.IsZero() {
		lastH = pet["last_decay_at"].Time()
	}
	if lastH.IsZero() {
		lastH = now.Time
	}
	lastA := pet["last_affection_at"].Time()
	if lastA.IsZero() {
		lastA = pet["last_decay_at"].Time()
	}
	if lastA.IsZero() {
		lastA = now.Time
	}

	hHours := now.Time.Sub(lastH).Hours()
	aHours := now.Time.Sub(lastA).Hours()
	if hHours < 1 && aHours < 1 {
		return
	}

	nh, na := hunger, aff
	var setLastH, setLastA bool
	var newLastH, newLastA time.Time
	if hunger > 0 {
		hPeriods := int(hHours) / petHungerHours
		if hPeriods > 0 {
			// 每 24h 扣当前饱食度的 20%（指数衰减，多周期复利）
			decay := int(float64(hunger) * (1 - math.Pow(1-petHungerPct, float64(hPeriods))))
			nh = hunger - decay
			if nh < 0 {
				nh = 0
			}
			// 仅推进已结算的整周期，保留零头避免时间漂移
			newLastH = lastH.Add(time.Duration(hPeriods*petHungerHours) * time.Hour)
			setLastH = true
		}
	}
	if aff > 0 {
		aPeriods := int(aHours) / petAffectionHours
		if aPeriods > 0 {
			na = aff - aPeriods
			if na < 0 {
				na = 0
			}
			newLastA = lastA.Add(time.Duration(aPeriods*petAffectionHours) * time.Hour)
			setLastA = true
		}
	}

	data := g.Map{}
	if nh != hunger {
		data["hunger"] = nh
	}
	if na != aff {
		data["affection"] = na
	}
	if setLastH {
		data["last_hunger_at"] = newLastH
	}
	if setLastA {
		data["last_affection_at"] = newLastA
	}
	if len(data) == 0 {
		return
	}
	if _, err := g.DB().Model("pets").Ctx(ctx).Where("child_id", childID).Data(data).Update(); err != nil {
		gLog("宠物衰减失败: %v", err)
		return
	}
	pet["hunger"] = g.NewVar(nh)
	pet["affection"] = g.NewVar(na)
	if setLastH {
		pet["last_hunger_at"] = g.NewVar(newLastH)
	}
	if setLastA {
		pet["last_affection_at"] = g.NewVar(newLastA)
	}

	// 惩罚：饱食度首次降到0 → 清空积分
	if nh == 0 && hunger > 0 {
		s.clearPoints(ctx, childID)
	}
	// 惩罚：好感度首次降到0 → 清空星星
	if na == 0 && aff > 0 {
		s.clearStars(ctx, childID)
	}
}

// clearPoints 清空学生积分（删除全部积分流水）。
func (s *sStudyPlanet) clearPoints(ctx context.Context, childID int) {
	if _, err := daoPointsLog.Ctx(ctx).Where("child_id", childID).Delete(); err != nil {
		gLog("宠物惩罚-清空积分失败 child=%d: %v", childID, err)
		return
	}
	gLog("宠物惩罚：饱食度归零，已清空 child=%d 的积分", childID)
}

// clearStars 清空学生学习星星（关卡星级归零）。
func (s *sStudyPlanet) clearStars(ctx context.Context, childID int) {
	if _, err := daoSessions.Ctx(ctx).Where("child_id", childID).Data(g.Map{"stars": 0, "max_combo": 0}).Update(); err != nil {
		gLog("宠物惩罚-清空星星失败 child=%d: %v", childID, err)
		return
	}
	gLog("宠物惩罚：好感度归零，已清空 child=%d 的星星", childID)
}

// petToRes 宠物行 → 响应结构。
func (s *sStudyPlanet) petToRes(pet gdb.Record) *v1.Pet {
	spCode := pet["species"].String()
	var sp petSpecies
	for _, x := range petSpeciesList {
		if x.Code == spCode {
			sp = x
			break
		}
	}
	if sp.Code == "" {
		sp = petSpeciesList[0]
	}
	hunger := pet["hunger"].Int()
	aff := pet["affection"].Int()
	mood, text := petMood(hunger, aff)
	lastFed := ""
	if t := pet["last_fed_at"].Time(); !t.IsZero() && t.Year() > 2000 {
		lastFed = t.Format("01-02 15:04")
	}
	foodInv := map[string]int{
		"apple":  pet["food_apple"].Int(),
		"fish":   pet["food_fish"].Int(),
		"milk":   pet["food_milk"].Int(),
		"star":   pet["food_star"].Int(),
		"cake":   pet["food_cake"].Int(),
		"hotpot": pet["food_hotpot"].Int(),
	}
	return &v1.Pet{
		ChildID: pet["child_id"].Int(), Name: pet["name"].String(),
		Species: sp.Code, SpeciesName: sp.Name, Emoji: sp.Emoji,
		Level: pet["level"].Int(), Exp: pet["exp"].Int(),
		ExpMax: petExpNeed(pet["level"].Int()),
		ExpToNext: petExpNeed(pet["level"].Int()) - pet["exp"].Int(),
		Hunger: hunger, Affection: aff,
		Mood: mood, MoodText: text,
		FedCount: pet["fed_count"].Int(), LastFedAt: lastFed,
		FoodInv: foodInv,
	}
}

// petLevelUp 通用升级结算：exp 累计升级（Lv 上限 99），返回新经验/新等级/是否升级。
func petLevelUp(curExp, curLevel, delta int) (newExp, newLevel int, levelUp bool) {
	newExp, newLevel = curExp+delta, curLevel
	for newLevel < 99 && newExp >= petExpNeed(newLevel) {
		newExp -= petExpNeed(newLevel)
		newLevel++
		levelUp = true
	}
	return
}

// spEmojiOf 种类 code → emoji（找不到给默认）。
func spEmojiOf(code string) string {
	for _, x := range petSpeciesList {
		if x.Code == code {
			return x.Emoji
		}
	}
	return "🐱"
}