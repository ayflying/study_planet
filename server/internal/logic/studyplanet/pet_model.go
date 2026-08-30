// Package studyplanet 宠物核心模型：种类/食物表、等级公式、心情推导、读写与惰性衰减。
package studyplanet

import (
	"context"
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

// petExpNeed 升到 level+1 所需累计经验：20 + (level-1)*15。
func petExpNeed(level int) int { return 20 + (level-1)*15 }

// petHungerRate 每小时饱食度衰减。
const petHungerRate = 4

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
		"child_id":      childID,
		"name":          sp.Name,
		"species":       sp.Code,
		"hunger":        80,
		"affection":     20,
		"last_decay_at": gtime.Now(),
	}).Insert(); err != nil {
		return nil, err
	}
	return g.DB().Model("pets").Ctx(ctx).Where("child_id", childID).One()
}

// petTick 惰性结算：按距上次的时长衰减饱食度（下限 0），并更新衰减时间。
func (s *sStudyPlanet) petTick(ctx context.Context, pet gdb.Record) {
	hunger := pet["hunger"].Int()
	last := pet["last_decay_at"].Time()
	hours := time.Since(last).Hours()
	if hours < 1 || hunger <= 0 {
		return
	}
	decay := int(hours) * petHungerRate
	if decay <= 0 {
		return
	}
	nh := hunger - decay
	if nh < 0 {
		nh = 0
	}
	if _, err := g.DB().Model("pets").Ctx(ctx).Where("child_id", pet["child_id"].Int()).Data(g.Map{
		"hunger": nh, "last_decay_at": gtime.Now(),
	}).Update(); err != nil {
		gLog("宠物衰减失败: %v", err)
		return
	}
	pet["hunger"] = g.NewVar(nh)
	pet["last_decay_at"] = g.NewVar(gtime.Now())
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
	return &v1.Pet{
		ChildID: pet["child_id"].Int(), Name: pet["name"].String(),
		Species: sp.Code, SpeciesName: sp.Name, Emoji: sp.Emoji,
		Level: pet["level"].Int(), Exp: pet["exp"].Int(),
		ExpMax: petExpNeed(pet["level"].Int()),
		Hunger: hunger, Affection: aff,
		Mood: mood, MoodText: text,
		FedCount: pet["fed_count"].Int(), LastFedAt: lastFed,
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
