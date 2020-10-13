package coupon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-promotion/business/account"
	"github.com/gingerxman/ginger-promotion/business/coupon/params"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
	"github.com/gingerxman/gorm"
)

type CouponService struct {
	eel.ServiceBase
}

func NewCouponService(ctx context.Context) *CouponService {
	service := new(CouponService)
	service.Ctx = ctx
	return service
}

func (this *CouponService) DeliverOne(rule *Rule, userId int, sourceType string) (*Coupon, error) {
	coupons, err := this.Deliver(rule, []int{userId}, 1, "initiative")
	if err != nil {
		return nil, err
	}
	
	return coupons[0], nil
}

// 发放优惠券
func (this *CouponService) Deliver(rule *Rule, userIds []int, countPerUser int, sourceType string) ([]*Coupon, error) {
	// 检查优惠券
	if rule == nil{
		return nil, eel.NewBusinessError("coupon:invalid_rule", fmt.Sprintf("错误的优惠券"))
	}
	if !rule.IsActive{
		return nil, eel.NewBusinessError("coupon:rule_deleted", fmt.Sprintf("优惠券已下架"))
	}
	if rule.IsDeleted {
		return nil, eel.NewBusinessError("coupon:rule_invalid", fmt.Sprintf("优惠券不存在"))
	}
	if rule.ReceiveStatus == "not_start" {
		return nil, eel.NewBusinessError("coupon:receive_not_start", fmt.Sprintf("未到领取时间"))
	}
	if rule.ReceiveStatus == "end" {
		return nil, eel.NewBusinessError("coupon:receive_end", fmt.Sprintf("已过领取时间"))
	}

	totalNeedCount := int64(len(userIds) * countPerUser)
	if rule.RemainCouponCount < totalNeedCount{
		return nil, eel.NewBusinessError("coupon:not_enough", fmt.Sprintf("库存不足"))
	}

	// 限领
	conditions := rule.Conditions
	type2condition := make(map[string]*simplejson.Json)
	for _, condition := range conditions{
		type2condition[condition.Type] = condition.Data
	}
	limitCountPerUser, ok := type2condition["count_per_user"]
	if ok{
		limitInfo := limitCountPerUser.Get("count")
		if limitInfo != nil{
			dataStr,_ := json.Marshal(limitCountPerUser)
			var aa struct {
				Count json.Number `json:"count"`
			}
			_ = json.Unmarshal([]byte(string(dataStr)), &aa)
			limitCount, _ := aa.Count.Int64()
			// 获取当前用户领取的张数
			o := eel.GetOrmFromContext(this.Ctx)
			var models []*m_coupon.Coupon
			db := o.Model(&m_coupon.Coupon{}).Where("user_id__in",userIds).Where("rule_id", rule.Id).Find(&models)
			if db.Error != nil {
				eel.Logger.Error(db.Error)
				panic(eel.NewBusinessError("coupon:coupons_failed", fmt.Sprintf("获取该用户领取过的优惠券失败")))
			}
			userId2couponCount := make(map[int]int)
			for _, coupon := range models {
				count, ok := userId2couponCount[coupon.UserId]
				if ok {
					userId2couponCount[coupon.UserId] = count + 1
				} else {
					userId2couponCount[coupon.UserId] = 1
				}
			}
			for _, userId := range userIds {
				couponCount, ok := userId2couponCount[userId]
				if !ok {
					couponCount = 0
				}
				if int(couponCount + countPerUser) > int(limitCount) {
					return nil, eel.NewBusinessError("coupon:exceed_limit_count", fmt.Sprintf("超过了可领取的次数"))
				}
			}
		}
	}

	// 发放优惠券
	coupons := make([]*Coupon, 0)
	for _, userId := range userIds{
		// 生成优惠券
		gnerateCoupons, err := rule.GenerateCoupons(userId, countPerUser, sourceType)
		if err != nil{
			return nil, err
		}
		for _, coupon := range gnerateCoupons{
			coupons = append(coupons, coupon)
		}
	}
	return coupons, nil
}

//恢复优惠券
func (this *CouponService) RestoredCoupon(user *account.User, coupon *Coupon, ) error {
	o := eel.GetOrmFromContext(this.Ctx)
	db:= o.Model(&m_coupon.Coupon{}).Where("user_id", user.GetId()).Where("id", coupon.Id).Update(gorm.Params{
		"status": m_coupon.COUPON_STATUS_UNUSED,
		"order_bid": "",
		"consumed_at": nil,
		"deduction_money": 0,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		bErr := eel.NewBusinessError("coupon:restored_failed", fmt.Sprintf("恢复优惠券失败"))
		return bErr
	}
	// TODO: 恢复account账户余额
	//isSuccess := common.NewPlutusService(this.Ctx).DoRestoredTransfer(
	//	user.PlatformId, fmt.Sprintf("cash.coupon.%d", coupon.RuleId), user.GetId(), coupon.DeductionMoney, coupon.OrderBid, "使用优惠券下单失败后的转账")
	//if !isSuccess {
	//	eel.Logger.Error(err)
	//	panic(eel.NewBusinessError("rule:create_fail", fmt.Sprintf("优惠券规则更新账户失败")))
	//}
	return nil
}


func (this *CouponService) DiscardCouponsByRules(ruleIds []int) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.Coupon{}).Where("rule_id__in", ruleIds).Update(gorm.Params{
		"status": m_coupon.COUPON_STATUS_DISCARD,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

// 计算可抵扣的金额
func (this *CouponService) Calculate(user *account.User, coupon *Coupon, couponInfo *params.CouponInfo, orderInfo *params.OrderInfo, products []*params.ProductInfo) (float64, error) {
	//isValidate, err := this.IsValidate(user, coupon, couponInfo, orderInfo, products)
	//if isValidate{
	//	deductionMoney
	//	return deductionMoney, nil
	//} else {
	//	return 0, err
	//}
	return 0, nil
}


func init() {
}
