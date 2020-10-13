package coupon

import (
	"context"
	"github.com/gingerxman/eel"

	
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
)

type FillCouponService struct {
	eel.ServiceBase
}

func NewFillCouponService(ctx context.Context) *FillCouponService {
	service := new(FillCouponService)
	service.Ctx = ctx
	return service
}

func (this *FillCouponService) Fill(coupons []*Coupon, option eel.FillOption) {
	if len(coupons) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, coupon := range coupons {
		ids = append(ids, coupon.Id)
	}

	if enableOption, ok := option["with_rule"]; ok && enableOption {
		this.fillRule(coupons, ids)
	}
	return
}


func (this *FillCouponService) fillRule(coupons []*Coupon, ids []int) {
	//获取关联的id集合
	ruleIds := make([]int, 0)
	for _, coupon := range coupons {
		ruleIds = append(ruleIds, coupon.RuleId)
	}

	//从db中获取数据集合
	var models []*m_coupon.Rule
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.Rule{}).Where("id__in", ruleIds).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}

	//构建<id, model>
	id2model := make(map[int]*m_coupon.Rule)
	for _, model := range models {
		id2model[model.Id] = model
	}

	//填充coupon的Rule对象
	rules := make([]*Rule, 0)
	for _, coupon := range coupons {
		if model, ok := id2model[coupon.RuleId]; ok {
			rule := NewRuleFromModel(this.Ctx, model)
			coupon.Rule = rule
			rules = append(rules, rule)
		}
	}
	
	NewFillRuleService(this.Ctx).Fill(rules, eel.FillOption{
		"with_condition": true,
	})
}


func init() {
}
