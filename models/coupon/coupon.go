package coupon

import (
	"time"
	"github.com/gingerxman/eel"
)


//Rule Model
type Rule struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	Name string 
	Count int 
	Desc string 
	Remark string 
	IsActive bool `gorm:"default:true"`
	IsDeleted bool `gorm:"default:false"`
	StartDate time.Time `gorm:"type:datetime"`
	EndDate time.Time `gorm:"type:datetime"`
	ReceiveStartDate time.Time `gorm:"type:datetime"`
	ReceiveEndDate time.Time `gorm:"type:datetime"`
}
func (self *Rule) TableName() string {
	return "coupon_rule"
}


// 优惠券规则条件类型
const CONDITION_TYPE_SINGLE_PRODUCT = 1 // 单品
const CONDITION_TYPE_DISCOUNT = 2 // 折扣
const CONDITION_TYPE_REACH_DEDUCTION = 3 // 满减（满xx后抵扣）
const CONDITION_TYPE_DEDUCTION = 4 // 抵扣
const CONDITION_TYPE_COUNT_PER_USER = 5 // 每人限领张数

var CONDITION_TYPE_CODE2TYPE = map[string]int{
	"single_product": CONDITION_TYPE_SINGLE_PRODUCT,
	"discount": CONDITION_TYPE_DISCOUNT,
	"reach_deduction": CONDITION_TYPE_REACH_DEDUCTION,
	"deduction": CONDITION_TYPE_DEDUCTION,
	"count_per_user": CONDITION_TYPE_COUNT_PER_USER,
}
var CONDITION_TYPE2CODE = map[int]string {
	CONDITION_TYPE_SINGLE_PRODUCT: "single_product",
	CONDITION_TYPE_DISCOUNT: "discount",
	CONDITION_TYPE_REACH_DEDUCTION: "reach_deduction",
	CONDITION_TYPE_DEDUCTION: "deduction",
	CONDITION_TYPE_COUNT_PER_USER: "count_per_user",
}
//RuleCondition Model
type RuleCondition struct {
	eel.Model
	RuleId int `gorm:index`
	Type int `gorm:index`
	Data string
}
func (self *RuleCondition) TableName() string {
	return "coupon_rule_condition"
}


// 优惠券状态
//const COUPON_STATUS_UNRECEIVED = 0 // unreceived未领取
const COUPON_STATUS_UNUSED = 0 // 已领取
const COUPON_STATUS_USED = 1 // 已使用
const COUPON_STATUS_EXPIRED = 2 // 已过期
const COUPON_STATUS_INVALID = 3 // 已失效
const COUPON_STATUS_DISCARD = 4 // 作废
var COUPON_STATUS_CODE2STATUS = map[string]int{
	"unused": COUPON_STATUS_UNUSED,
	"used": COUPON_STATUS_USED,
	"expired": COUPON_STATUS_EXPIRED,
	"invalid": COUPON_STATUS_INVALID,
	"discard": COUPON_STATUS_DISCARD,
}
var COUPON_STATUS2CODE = map[int]string {
	COUPON_STATUS_UNUSED: "unused",
	COUPON_STATUS_USED: "used",
	COUPON_STATUS_EXPIRED: "expired",
	COUPON_STATUS_INVALID: "invalid",
	COUPON_STATUS_DISCARD: "discard",
}
var COUPON_STATUS2CODE2STR = map[string]string {
	"unused": "未使用",
	"used": "已使用",
	"expired": "已过期",
	"invalid": "已失效",
	"discard": "作废",
}
// 优惠券的来源
const COUPON_SOURCE_TYPE_INITIATIVE = 1 // 主动领取的
const COUPON_SOURCE_TYPE_ASSIGNED = 2 // 分配
var COUPON_SOURCE_TYPE_CODE2TYPE = map[string]int{
	"initiative": COUPON_SOURCE_TYPE_INITIATIVE,
	"assigned": COUPON_SOURCE_TYPE_ASSIGNED,
}
var COUPON_SOURCE_TYPE2CODE = map[int]string {
	COUPON_SOURCE_TYPE_INITIATIVE: "initiative",
	COUPON_SOURCE_TYPE_ASSIGNED: "assigned",
}
//Coupon Model
type Coupon struct {
	eel.Model
	RuleId int `gorm:index`
	UserId int             `gorm:"default:0;index"`
	SourceType int         `gorm:"default:0"`
	Code string            `gorm:"size:50;index"`
	Status int             `gorm:"default:0"`
	OrderBid string        `gorm:"size:32;default:'';index"`
	DeductionMoney float64 `sql:"type:float(255,2);"`
	ConsumedAt time.Time   `gorm:"type:datetime"`
}
func (self *Coupon) TableName() string {
	return "coupon_coupon"
}


func init() {
	eel.RegisterModel(new(Rule))
	eel.RegisterModel(new(RuleCondition))
	eel.RegisterModel(new(Coupon))
	
}
