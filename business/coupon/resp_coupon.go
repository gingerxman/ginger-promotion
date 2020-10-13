package coupon

import (
	"github.com/bitly/go-simplejson"
)

type RRule struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Count int `json:"count"`
	Desc string `json:"desc"`
	Remark string `json:"remark"`
	IsActive bool `json:"is_active"`
	IsDeleted bool `json:"is_deleted"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	CreatedAt string `json:"created_at"`
	Status string `json:"status"`
	Conditions []*RRuleCondition `json:"conditions"`
	RemainCouponCount int64 `json:"remain_coupon_count"`
	ReceiveStartDate string `json:"receive_start_date"`
	ReceiveEndDate string `json:"receive_end_date"`
	ReceiveStatus string `json:"receive_status"`
	CouponLinkUrl string `json:"coupon_link_url"`
}

type RRuleCondition struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Data *simplejson.Json `json:"data"`
	CreatedAt string `json:"created_at"`
}

type RCoupon struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	SourceType string `json:"source_type"`
	Code string `json:"code"`
	Status string `json:"status"`
	OrderBid string `json:"order_bid"`
	DeductionMoney float64 `json:"deduction_money"`
	ConsumedAt string `json:"consumed_at"`
	Rule *RRule `json:"rule"`
	CreatedAt string `json:"created_at"`
}


func init() {
}
