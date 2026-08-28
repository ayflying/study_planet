package v1

import "github.com/gogf/gf/v2/frame/g"

// Reward 奖励项。
type Reward struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CostPoints int    `json:"cost_points"`
	Status     string `json:"status"` // active | redeemed
}

// ListRewardsReq 奖励列表（公开接口）。
type ListRewardsReq struct {
	g.Meta `path:"/rewards" method:"get" tags:"Reward" summary:"奖励列表"`
}
type ListRewardsRes []Reward

// AddRewardReq 添加奖励（家长鉴权）。
type AddRewardReq struct {
	g.Meta     `path:"/parent/rewards" method:"post" tags:"Reward" summary:"添加奖励"`
	Name       string `json:"name" v:"required"`
	CostPoints int    `json:"cost_points"`
}
type AddRewardRes struct {
	OK bool `json:"ok"`
}

// RedeemReq 学生兑换奖励（公开接口，需家长确认）。
type RedeemReq struct {
	g.Meta `path:"/rewards/:id/redeem" method:"post" tags:"Reward" summary:"兑换奖励"`
	ID     int `in:"path" json:"-"`
}
type RedeemRes struct {
	OK      bool   `json:"ok"`
	Pending bool   `json:"pending"`
	Message string `json:"message"`
}

// ConfirmRedemptionReq 家长确认兑换（家长鉴权）。
type ConfirmRedemptionReq struct {
	g.Meta `path:"/parent/redemptions/:id/confirm" method:"post" tags:"Reward" summary:"确认兑换"`
	ID     int `in:"path" json:"-"`
}
type ConfirmRedemptionRes struct {
	OK bool `json:"ok"`
}
