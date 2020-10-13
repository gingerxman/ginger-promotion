package coupon

import (
	"context"
	
	"github.com/gingerxman/eel"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
)

type FillRuleService struct {
	eel.ServiceBase
}

func NewFillRuleService(ctx context.Context) *FillRuleService {
	service := new(FillRuleService)
	service.Ctx = ctx
	return service
}

func (this *FillRuleService) Fill(rules []*Rule, option eel.FillOption) {
	if len(rules) == 0 {
		return
	}

	ids := make([]int, 0)
	for _, rule := range rules {
		ids = append(ids, rule.Id)
	}

	if enableOption, ok := option["with_condition"]; ok && enableOption {
		this.fillConditions(rules, ids)
	}
	return
}

func (this *FillRuleService) fillConditions(rules []*Rule, ids []int) {

	//从db中获取数据集合
	var models []*m_coupon.RuleCondition
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.RuleCondition{}).Where("rule_id__in", ids).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return
	}

	//构建<id, []>
	ruleId2models := make(map[int][]*m_coupon.RuleCondition)
	for _, model := range models {
		ruleId2models[model.RuleId] = append(ruleId2models[model.RuleId], model)
	}

	//填充Rule的RuleCondition对象
	for _, rule := range rules {
		if models, ok := ruleId2models[rule.Id]; ok {
			for _, model := range models{
				rule.AddRuleCondition(NewRuleConditionFromModel(this.Ctx, model))
			}

		}
	}
}

func init() {
}
