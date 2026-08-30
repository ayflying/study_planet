package v1

import "github.com/gogf/gf/v2/frame/g"

// Pet 宠物信息。
type Pet struct {
	ChildID     int            `json:"child_id"`
	Name        string         `json:"name"`
	Species     string         `json:"species"`
	SpeciesName string         `json:"species_name"`
	Emoji       string         `json:"emoji"`
	Level       int            `json:"level"`
	Exp         int            `json:"exp"`
	ExpMax      int            `json:"exp_max"`
	Hunger      int            `json:"hunger"`
	Affection   int            `json:"affection"`
	Mood        string         `json:"mood"`
	MoodText    string         `json:"mood_text"`
	FedCount    int            `json:"fed_count"`
	LastFedAt   string         `json:"last_fed_at"`
	FoodInv     map[string]int `json:"food_inventory"` // 食物库存，如 {"apple":1}
	SnackMsg    string         `json:"snack_msg,omitempty"` // 掉落提示/温馨提醒
}

// PetFood 投喂食物项。
type PetFood struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Emoji       string `json:"emoji"`
	Hunger      int    `json:"hunger"`
	Affection   int    `json:"affection"`
	Exp         int    `json:"exp"`
	Description string `json:"description"`
}

// PetGetReq 查看我的宠物（首次访问自动领取）。
type PetGetReq struct {
	g.Meta    `path:"/pet" method:"get" tags:"Pet" summary:"我的宠物"`
	StudentID int `json:"student_id" in:"query"`
}
type PetGetRes Pet

// PetFeedReq 投喂宠物：food 为食物 id。
type PetFeedReq struct {
	g.Meta    `path:"/pet/feed" method:"post" tags:"Pet" summary:"投喂宠物"`
	StudentID int    `json:"student_id" in:"query"`
	Food      string `json:"food" v:"required"`
}
type PetFeedRes struct {
	Pet        Pet    `json:"pet"`
	Message    string `json:"message"`
	LevelUp    bool   `json:"level_up"`
	FedBurst   bool   `json:"fed_burst"` // 好感度突破整数关口时触发特效（如升到 50/80/100）
}

// PetRenameReq 给宠物改名。
type PetRenameReq struct {
	g.Meta    `path:"/pet/rename" method:"post" tags:"Pet" summary:"宠物改名"`
	StudentID int    `json:"student_id" in:"query"`
	Name      string `json:"name" v:"required"`
}
type PetRenameRes Pet

// PetDecayReq 空闲衰减（打开宠物页时惰性结算饱食度/心情）。
type PetDecayReq struct {
	g.Meta    `path:"/pet/tick" method:"post" tags:"Pet" summary:"宠物状态结算"`
	StudentID int `json:"student_id" in:"query"`
}
type PetDecayRes Pet

// PetFoodsReq 可投喂的食物清单。
type PetFoodsReq struct {
	g.Meta `path:"/pet/foods" method:"get" tags:"Pet" summary:"食物清单"`
}
type PetFoodsRes []PetFood
