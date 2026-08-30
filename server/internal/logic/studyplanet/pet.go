// Package studyplanet 宠物接口：查看（自动领取+惰性衰减）、投喂、改名、经验联动、食物清单。
package studyplanet

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "studyplanet/api/studyplanet/v1"
)

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
	r.SnackMsg = "温馨提示：饱食度降为0时会清空积分，好感度降为0时会清空星星哦！"
	return (*v1.PetGetRes)(r), nil
}

// PetFeed 投喂（需消耗库存食物，完成任务/对战/学习可掉落获得）。
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

	// 检查库存
	col := inventoryColumn(food.ID)
	if col == "" {
		return nil, errParam("未知的食物类型")
	}
	qty := pet[col].Int()
	if qty <= 0 {
		return nil, errParam("这种零食吃完了，快去完成任务、对战或学习获取零食吧！")
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
	// 经验与升级（复用公共升级结算）
	newExp, newLevel, levelUp := petLevelUp(pet["exp"].Int(), pet["level"].Int(), food.Exp)
	fedBurst := aff < 100 && newAff == 100
	if _, err := g.DB().Model("pets").Ctx(ctx).Where("child_id", cid).Data(g.Map{
		"hunger": newHunger, "affection": newAff, "exp": newExp, "level": newLevel,
		"fed_count": pet["fed_count"].Int() + 1, "last_fed_at": gtime.Now(), "last_decay_at": gtime.Now(),
		col: qty - 1, // 消耗1个库存
	}).Update(); err != nil {
		return nil, gerror.Wrap(err, "投喂失败")
	}
	pet["hunger"], pet["affection"], pet["exp"], pet["level"], pet["fed_count"] = g.NewVar(newHunger), g.NewVar(newAff), g.NewVar(newExp), g.NewVar(newLevel), g.NewVar(pet["fed_count"].Int()+1)
	pet[col] = g.NewVar(qty - 1)
	msg := fmt.Sprintf("%s%s 吃得津津有味！", spEmojiOf(pet["species"].String()), pet["name"].String())
	if hunger >= 100 {
		msg = fmt.Sprintf("它已经饱了，但还是开心地收下了 %s%s！", food.Emoji, food.Name)
	}
	if levelUp {
		msg += fmt.Sprintf(" 🎉 升到 Lv.%d 了！", newLevel)
	}
	return &v1.PetFeedRes{Pet: *s.petToRes(pet), Message: msg, LevelUp: levelUp, FedBurst: fedBurst}, nil
}

// addSnackDrop 一次"零食掉落"：命中则给该学生宠物加1个对应零食，返回掉落的 foodID（空表示未掉落）。
func (s *sStudyPlanet) addSnackDrop(ctx context.Context, childID int) string {
	if childID <= 0 {
		return ""
	}
	food := rollSnack()
	if food == "" {
		return ""
	}
	col := inventoryColumn(food)
	if col == "" {
		return ""
	}
	// 直接 SQL 累加，避免竞态
	q := "UPDATE pets SET " + col + " = " + col + " + 1 WHERE child_id=?"
	if _, err := g.DB().Exec(ctx, q, childID); err != nil {
		gLog("零食掉落累加失败 child=%d food=%s: %v", childID, food, err)
		return ""
	}
	gLog("零食掉落：child=%d 获得 %s", childID, food)
	return food
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
	newExp, newLevel, _ := petLevelUp(pet["exp"].Int(), pet["level"].Int(), delta)
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

// PetFoods 食物清单（前端宠物页渲染用，公开端点）。
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

// ExternalSnackDrop 供 battle 引擎在胜利后给学生掉落零食。
func ExternalSnackDrop() func(childID int) string {
	return func(childID int) string {
		return Study().addSnackDrop(gctx.New(), childID)
	}
}