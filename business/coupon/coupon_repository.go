package coupon

import (
	"context"
	"github.com/gingerxman/ginger-promotion/business"
	"github.com/gingerxman/ginger-promotion/business/account"
	"github.com/gingerxman/ginger-promotion/business/product"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
	
	
	"github.com/gingerxman/eel"
)

type CouponRepository struct {
	eel.RepositoryBase
}

func NewCouponRepository(ctx context.Context) *CouponRepository {
	repository := new(CouponRepository)
	repository.Ctx = ctx
	return repository
}

func (this *CouponRepository) GetCoupons(filters eel.Map, orderExprs ...string) []*Coupon {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.Coupon{})
	
	var models []*m_coupon.Coupon
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	db = db.Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	coupons := make([]*Coupon, 0)
	for _, model := range models {
		coupons = append(coupons, NewCouponFromModel(this.Ctx, model))
	}
	return coupons
}

func (this *CouponRepository) GetPagedCoupons(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*Coupon, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.Coupon{})
	
	var models []*m_coupon.Coupon
	if len(filters) > 0 {
		db = db.Where(filters)
	}
	for _, expr := range orderExprs {
		db = db.Order(expr)
	}
	paginateResult, db := eel.Paginate(db, page, &models)
	err := db.Error
	if err != nil {
		eel.Logger.Error(err)
		return nil, paginateResult
	}
	
	Coupons := make([]*Coupon, 0)
	for _, model := range models {
		Coupons = append(Coupons, NewCouponFromModel(this.Ctx, model))
	}
	return Coupons, paginateResult
}

//GetEnabledCouponsForCorp 获得启用的Coupon对象集合
func (this *CouponRepository) GetEnabledCouponsForCorp(rule *Rule, page *eel.PageInfo, filters eel.Map) ([]*Coupon, eel.INextPageInfo) {
	filters["rule_id"] = rule.Id
	
	return this.GetPagedCoupons(filters, page, "-id")
}

//GetCouponsForUser 获得启用的Coupon对象集合
func (this *CouponRepository) GetCouponsForUser(user *account.User, page *eel.PageInfo, filters eel.Map) ([]*Coupon, eel.INextPageInfo) {
	filters["user_id"] = user.GetId()
	status, ok := filters["status"]
	if ok{
		filters["status"] = m_coupon.COUPON_STATUS_CODE2STATUS[status.(string)]
	}
	
	return this.GetPagedCoupons(filters, page, "-id")
}

//GetCouponsForUserAndProducts 获得启用的Coupon对象集合
func (this *CouponRepository) GetCouponsForUserAndProducts(user *account.User, products []*product.Product) []*Coupon {
	page := &eel.PageInfo{
		Page: 1,
		CountPerPage: 999,
	}
	
	filters := eel.Map{
		"status": "unused",
	}
	coupons, _ := this.GetCouponsForUser(user, page, filters)
	
	//根据pool product ids过滤coupon
	id2exist := make(map[int]bool)
	for _, product := range products {
		id2exist[product.Id] = true
		id2exist[product.SourceProductId] = true
	}
	
	NewFillCouponService(this.Ctx).Fill(coupons, eel.FillOption{
		"with_rule": true,
	})
	filteredCoupons := make([]*Coupon, 0)
	for _, coupon := range coupons {
		if coupon.IsForSingleProduct() {
			poolProductId := coupon.GetTargetPoolProductId()
			if _, ok := id2exist[poolProductId]; ok {
				filteredCoupons = append(filteredCoupons, coupon)
			}
		} else {
			filteredCoupons = append(filteredCoupons, coupon)
		}
	}
	
	return filteredCoupons
}

//GetAllCouponsForCorp 获得所有Coupon对象集合
func (this *CouponRepository) GetAllCouponsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*Coupon, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedCoupons(filters, page, "-id")
}

//GetCouponInCorp 根据id和corp获得Coupon对象
func (this *CouponRepository) GetCouponInCorp(corp business.ICorp, id int) *Coupon {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	Coupons := this.GetCoupons(filters)
	
	if len(Coupons) == 0 {
		return nil
	} else {
		return Coupons[0]
	}
}

//GetCoupon 根据id和corp获得Coupon对象
func (this *CouponRepository) GetCoupon(id int) *Coupon {
	filters := eel.Map{
		"id": id,
	}
	
	Coupons := this.GetCoupons(filters)
	
	if len(Coupons) == 0 {
		return nil
	} else {
		return Coupons[0]
	}
}

//GetCoupon 根据code获得Coupon对象
func (this *CouponRepository) GetCouponByCode(code string) *Coupon {
	filters := eel.Map{
		"code": code,
	}

	Coupons := this.GetCoupons(filters)

	if len(Coupons) == 0 {
		return nil
	} else {
		return Coupons[0]
	}
}

//GetCouponByCodeForUser 根据code获得Coupon对象
func (this *CouponRepository) GetCouponByCodeForUser(user business.IUser, code string) *Coupon {
	filters := eel.Map{
		"code": code,
		"user_id": user.GetId(),
	}

	Coupons := this.GetCoupons(filters)

	if len(Coupons) == 0 {
		return nil
	} else {
		return Coupons[0]
	}
}

func init() {
}
