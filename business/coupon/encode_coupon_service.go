package coupon

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodeCouponService struct {
	eel.ServiceBase
}

func NewEncodeCouponService(ctx context.Context) *EncodeCouponService {
	service := new(EncodeCouponService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeCouponService) Encode(coupon *Coupon) *RCoupon {
	if coupon == nil {
		return nil
	}
	rRule := NewEncodeRuleService(this.Ctx).Encode(coupon.Rule)
	var consumedAt string
	if coupon.ConsumedAt.IsZero() {
		consumedAt = ""
	} else {
		consumedAt = coupon.ConsumedAt.Format("2006-01-02 15:04:05")
	}
	
	status := coupon.Status
	if coupon.Rule.IsDeleted {
		status = "discard"
	}
	return &RCoupon{
		Id: coupon.Id,
		UserId: coupon.UserId,
		SourceType: coupon.SourceType,
		Code: coupon.Code,
		Status: status,
		OrderBid: coupon.OrderBid,
		DeductionMoney: coupon.DeductionMoney,
		ConsumedAt: consumedAt,
		Rule: rRule,
		CreatedAt: coupon.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeCouponService) EncodeMany(coupons []*Coupon) []*RCoupon {
	rDatas := make([]*RCoupon, 0)
	for _, coupon := range coupons {
		rDatas = append(rDatas, this.Encode(coupon))
	}
	
	return rDatas
}

func init() {
}
