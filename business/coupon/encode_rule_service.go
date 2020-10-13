package coupon

import (
	"context"
	"fmt"
	"github.com/gingerxman/eel"
)

type EncodeRuleService struct {
	eel.ServiceBase
}

func NewEncodeRuleService(ctx context.Context) *EncodeRuleService {
	service := new(EncodeRuleService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeRuleService) Encode(rule *Rule) *RRule {
	if rule == nil {
		return nil
	}

	conditionDatas := make([]*RRuleCondition, 0)
	if len(rule.Conditions) > 0 {
		for _, resource := range rule.Conditions {
			conditionDatas = append(conditionDatas, &RRuleCondition{
				Id:   resource.Id,
				Type: resource.Type,
				Data:  resource.Data,
			})
		}
	}

	return &RRule{
		Id: rule.Id,
		Name: rule.Name,
		Count: rule.Count,
		Desc: rule.Desc,
		Remark: rule.Remark,
		IsActive: rule.IsActive,
		IsDeleted: rule.IsDeleted,
		StartDate: rule.StartDate.Format("2006-01-02 15:05"),
		EndDate: rule.EndDate.Format("2006-01-02 15:05"),
		CreatedAt: rule.CreatedAt.Format("2006-01-02 15:05"),
		Status: rule.Status,
		Conditions: conditionDatas,
		RemainCouponCount: rule.RemainCouponCount,
		ReceiveStartDate: rule.ReceiveStartDate.Format("2006-01-02 15:05"),
		ReceiveEndDate: rule.ReceiveEndDate.Format("2006-01-02 15:05"),
		ReceiveStatus: rule.ReceiveStatus,
		CouponLinkUrl: fmt.Sprintf("http://%s/user_coupon/?rule_id=%d", "", rule.Id), // 域名暂时不知
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeRuleService) EncodeMany(rules []*Rule) []*RRule {
	rDatas := make([]*RRule, 0)
	for _, rule := range rules {
		rDatas = append(rDatas, this.Encode(rule))
	}
	
	return rDatas
}

func init() {
}
