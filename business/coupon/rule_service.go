package coupon

import (
	"context"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
	"fmt"
	
	"github.com/gingerxman/gorm"
	"github.com/gingerxman/eel"
)

type RuleService struct {
	eel.ServiceBase
}

func NewRuleService(ctx context.Context) *RuleService {
	service := new(RuleService)
	service.Ctx = ctx
	return service
}

//DeleteRule 根据id删除Rule对象
func (this *RuleService) DeleteRuleById(id int) bool {
	return this.DeleteRuleByIds([]int{id})
}

func (this *RuleService) DeleteRuleByIds(ids []int) bool {
	o := eel.GetOrmFromContext(this.Ctx)

	if len(ids) > 0 {
		db := o.Model(&m_coupon.Rule{}).Where("id__in", ids).Update(gorm.Params{
			"is_deleted": true,
		})
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("rule:update_fail", fmt.Sprintf("删除优惠券规则失败")))
		}
	}
	
	eel.Logger.Error(ids)
	NewCouponService(this.Ctx).DiscardCouponsByRules(ids)
	return true
}

func init() {
}
