// Package studyplanet 宠物模式：养成 + 投喂 + 好感度。
// 规则：
//   - 每个学生一只宠物（首次访问自动领取）；
//   - 饱食度随时间衰减（每小时 -4，下限 0），投喂回复；
//   - 好感度由投喂与闯关/对战表现累积（0-100）；
//   - 宠物经验主要来自主人答题：每答对一题 addPetExp +2（content/practice 链路调用）；
//   - 等级 = f(exp)：升级所需 exp 递增（Lv1→2 需 20，之后每级 +15）。
package studyplanet

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	v1 "studyplanet/api/studyplanet/v1"
)

// petSpecies 宠物种类表。
type petSpecies struct {
	Code string
	Name string
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
		"child_id": childID,
		"name":     sp.Name,
		"species":  sp.Code,
		"hunger":   80,
		"affection": 20,
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

// PetGet 查看宠物（自动领取 + 惰性衰减）。
func (s *sStudyPlanet) PetGet(ctx context.Context, req *v1.PetGetReq) (res *v1.PetGetRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid <= 0 {
		return nil, errNotFound("学生不存在")
	}
	pet, err := s.petOf(ctx, cid)
	if err != nil {
		return nil, gerror.Wrap(err, "查询宠物失败")
	}
	s.petTick(ctx, pet)
	r := s.petToRes(pet)
	return (*v1.PetGetRes)(r), nil
}

// PetFeed 投喂。
func (s *sStudyPlanet) PetFeed(ctx context.Context, req *v1.PetFeedReq) (res *v1.PetFeedRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid <= 0 {
		return nil, errNotFound("学生不存在")
	}
	pet, err := s.petOf(ctx, cid)
	if err != nil {
		return nil, gerror.Wrap(err, "查询宠物失败")
	}
	s.petTick(ctx, pet)
	var food petFood
	for _, f := range petFoodList {
		if f.ID == strings.TrimSpace(req.Food) {
			food = f
			break
		}
	}
	if food.ID == "" {
		return nil, errParam("没有这种食物哦")
	}
	hunger := pet["hunger"].Int()
	aff := pet["affection"].Int()
	// 已吃很饱时不再涨饱食度，但仍给少量好感（宠物会开心被投喂）
	if hunger >= 100 {
		food.Hunger = 0
	}
	newHunger := hunger + food.Hunger
	if newHunger > 100 {
		newHunger = 100
	}
	newAff := aff + food.Affection
	if newAff > 100 {
		newAff = 100
	}
	newExp := pet["exp"].Int() + food.Exp
	newLevel := pet["level"].Int()
	levelUp := false
	for newLevel < 99 && newExp >= petExpNeed(newLevel) {
		newExp -= petExpNeed(newLevel)
		newLevel++
		levelUp = true
	}
	fedBurst := aff < 100 && newAff == 100
	if _, err := g.DB().Model("pets").Ctx(ctx).Where("child_id", cid).Data(g.Map{
		"hunger": newHunger, "affection": newAff, "exp": newExp, "level": newLevel,
		"fed_count": pet["fed_count"].Int() + 1, "last_fed_at": gtime.Now(), "last_decay_at": gtime.Now(),
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "投喂失败")
	}
	pet["hunger"], pet["affection"], pet["exp"], pet["level"], pet["fed_count"] = g.NewVar(newHunger), g.NewVar(newAff), g.NewVar(newExp), g.NewVar(newLevel), g.NewVar(pet["fed_count"].Int()+1)
	msg := fmt.Sprintf("%s%s 吃得津津有味！", spEmojiOf(pet["species"].String()), pet["name"].String())
	if hunger >= 100 {
		msg = fmt.Sprintf("它已经饱了，但还是开心地收下了 %s%s！", food.Emoji, food.Name)
	}
	if levelUp {
		msg += fmt.Sprintf(" 🎉 升到 Lv.%d 了！", newLevel)
	}
	return &v1.PetFeedRes{Pet: *s.petToRes(pet), Message: msg, LevelUp: levelUp, FedBurst: fedBurst}, nil
}

// PetRename 改名。
func (s *sStudyPlanet) PetRename(ctx context.Context, req *v1.PetRenameReq) (res *v1.PetRenameRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || len([]rune(name)) > 12 {
		return nil, errParam("名字要 1-12 个字哦")
	}
	if _, err := g.DB().Model("pets").Ctx(ctx).Where("child_id", cid).Data(g.Map{"name": name}).Update(); err != nil {
		return nil, gerror.Wrap(err, "改名失败")
	}
	pet, err := s.petOf(ctx, cid)
	if err != nil {
		return nil, gerror.Wrap(err, "查询宠物失败")
	}
	r := s.petToRes(pet)
	return (*v1.PetRenameRes)(r), nil
}

// PetTick 打开宠物页时惰性结算（与 PetGet 同逻辑，供前端主动刷新）。
func (s *sStudyPlanet) PetTick(ctx context.Context, req *v1.PetDecayReq) (res *v1.PetDecayRes, err error) {
	cid, err := s.resolveChild(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if cid <= 0 {
		return nil, errNotFound("学生不存在")
	}
	pet, err := s.petOf(ctx, cid)
	if err != nil {
		return nil, gerror.Wrap(err, "查询宠物失败")
	}
	s.petTick(ctx, pet)
	r := s.petToRes(pet)
	return (*v1.PetDecayRes)(r), nil
}

// addPetExp 主人答题表现转化为宠物经验（content/practice/battle 链路调用）。
// correct=true 加 2 分经验，false 加 1（安慰奖）。
func (s *sStudyPlanet) addPetExp(ctx context.Context, childID int, correct bool) {
	if childID <= 0 {
		return
	}
	delta := 1
	if correct {
		delta = 2
	}
	pet, err := s.petOf(ctx, childID)
	if err != nil {
		return // 非关键路径
	}
	newExp := pet["exp"].Int() + delta
	newLevel := pet["level"].Int()
	for newLevel < 99 && newExp >= petExpNeed(newLevel) {
		newExp -= petExpNeed(newLevel)
		newLevel++
	}
	// 主人认真答题，宠物好感度也微涨（有上限）
	newAff := pet["affection"].Int()
	if correct && newAff < 100 {
		newAff++
	}
	if _, err := g.DB().Model("pets").Ctx(ctx).Where("child_id", childID).Data(g.Map{
		"exp": newExp, "level": newLevel, "affection": newAff,
	}).Update(); err != nil {
		gLog("宠物经验更新失败: %v", err)
	}
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

// petFoodResponse 食物清单（前端宠物页渲染用，公开端点 PetFoods）。
func (s *sStudyPlanet) PetFoods(ctx context.Context, req *v1.PetFoodsReq) (res *v1.PetFoodsRes, err error) {
	out := make(v1.PetFoodsRes, 0, len(petFoodList))
	for _, f := range petFoodList {
		out = append(out, v1.PetFood{
			ID: f.ID, Name: f.Name, Emoji: f.Emoji,
			Hunger: f.Hunger, Affection: f.Affection, Exp: f.Exp, Description: f.Desc,
		})
	}
	return &out, nil
}

