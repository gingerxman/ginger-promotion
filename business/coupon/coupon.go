package coupon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gingerxman/eel/snowflake"
	"github.com/gingerxman/ginger-promotion/business"
	"github.com/gingerxman/ginger-promotion/business/account"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
	"time"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/gorm"
)

var snowflakeNode, _ = snowflake.NewNode(1)

type Coupon struct {
	eel.EntityBase
	Id int
	UserId int
	SourceType string
	Code string
	Status string
	OrderBid string
	ConsumedAt time.Time
	DeductionMoney float64
	CreatedAt time.Time

	//foreign key
	RuleId int //refer to rule
	Rule *Rule
	isValidityChecked bool
	isValidate bool
	user *account.User
}

//Update 更新对象
func (this *Coupon) Update(
	userId int,
	sourceType int,
	code string,
	status int,
	orderBid string,
	consumedAt string,
	ruleId int,
) error {
	var model m_coupon.Coupon
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"user_id": userId,
		"source_type": sourceType,
		"code": code,
		"status": status,
		"order_bid": orderBid,
		"consumed_at": consumedAt,
		"rule_id": ruleId,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return errors.New("coupon:update_fail")
	}

	return nil
}

func (this *Coupon) UpdateStatus(status string) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_coupon.Coupon{}).Where(eel.Map{
		"id": this.Id,
	}).Update(gorm.Params{
		"status": m_coupon.COUPON_STATUS_CODE2STATUS[status],
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("coupon:update_status_fail", fmt.Sprintf("更新状态失败")))
	}
}

func (this *Coupon) GetRule() *Rule {
	if this.Rule == nil {
		repository := NewRuleRepository(this.Ctx)
		rule := repository.GetRule(this.RuleId)
		NewFillRuleService(this.Ctx).Fill([]*Rule{rule}, eel.FillOption{
			"with_condition": true,
		})
		this.Rule = rule
	}
	return this.Rule
}

func (this *Coupon) IsForSingleProduct() bool {
	rule := this.Rule
	for _, condition := range rule.Conditions {
		if condition.Type == "single_product" {
			return true
		}
	}
	
	return false
}

func (this *Coupon) GetTargetPoolProductId() int {
	rule := this.Rule
	for _, condition := range rule.Conditions {
		if condition.Type == "single_product" {
			poolProductIds := condition.Data.Get("products").MustArray()
			if len(poolProductIds) > 0 {
				id, _ := poolProductIds[0].(json.Number).Int64()
				return int(id)
			}
		}
	}
	
	return 0
}

// GetStatus 状态
func (this *Coupon,) GetStatus(model *m_coupon.Coupon) string{
	now := time.Now()
	endDate := this.GetRule().EndDate
	status := model.Status
	if endDate.Before(now) {
		o := eel.GetOrmFromContext(this.Ctx)
		db := o.Model(&m_coupon.Coupon{}).Where(eel.Map{
			"id": this.Id,
		}).Update(gorm.Params{
			"status": m_coupon.COUPON_STATUS_CODE2STATUS["expired"],
		})
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			panic(eel.NewBusinessError("coupon:update_status_fail", fmt.Sprintf("更新状态失败")))
		}
		status = m_coupon.COUPON_STATUS_EXPIRED
	}
	return m_coupon.COUPON_STATUS2CODE[status]
}

func (this *Coupon) CheckValidity(user *account.User, poolProductIds []int) (bool, error) {
	// 此优惠券是否是该使用者的
	this.isValidityChecked = true
	this.user = user
	userId := user.GetId()
	if this.UserId != userId {
		eel.Logger.Error("[validation] coupon(%s) not belong to current user({})", this.Code, userId)
		return false, eel.NewBusinessError("coupon:not_current_user", fmt.Sprintf("此优惠券不属于该用户"))
	}
	
	// 此优惠券是否是未使用的
	if this.Status != m_coupon.COUPON_STATUS2CODE[m_coupon.COUPON_STATUS_UNUSED] {
		eel.Logger.Error("[validation] coupon({}) is unusable", this.Code)
		return false, eel.NewBusinessError("coupon:unusable", fmt.Sprintf("此优惠券%s", m_coupon.COUPON_STATUS2CODE2STR[this.Status]))
	}
	
	conditions := this.GetRule().GetConditions()
	type2condition := make(map[string]*simplejson.Json)
	for _, condition := range conditions{
		type2condition[condition.Type] = condition.Data
	}
	// 此券是否是指定的商品可用的
	limitProductInfo, ok := type2condition["single_product"]
	if ok{
		limitProducts := limitProductInfo.Get("products")
		if limitProducts != nil {
			isValidate := false
			for _, poolProductId := range poolProductIds {
				for _, limitProductId := range limitProducts.MustArray() {
					limitProductId, _ := limitProductId.(json.Number).Int64()
					if int(limitProductId) == poolProductId {
						isValidate = true
					}
				}
			}
			if !isValidate {
				eel.Logger.Error("[validation] no product can use coupon({})", this.Code)
				return false, eel.NewBusinessError("coupon:product_not_use", fmt.Sprintf("该商品不可用此券"))
			}
		}
	}
	
	this.isValidate = true
	return true, nil
}

func (this *Coupon) UseByOrder(order business.IOrder) error {
	if !this.isValidityChecked {
		panic("check coupon validity first")
	}
	
	o := eel.GetOrmFromContext(this.Ctx)
	if !this.isValidate {
		return errors.New("coupon is not valid")
	}
	
	bid := order.GetBid()
	deductionMoney, err := this.GetDeductionMoney(order.GetDeductableMoney())
	if err != nil {
		eel.Logger.Error(err)
		bErr := eel.NewBusinessError("coupon:get_deduction_money_fail", fmt.Sprintf("获取抵扣金额失败"))
		return bErr
	}
	
	db := o.Model(&m_coupon.Coupon{}).Where("id", this.Id).Update(gorm.Params{
		"status": m_coupon.COUPON_STATUS_USED,
		"order_bid": bid,
		"consumed_at": time.Now(),
		"deduction_money": deductionMoney,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		bErr := eel.NewBusinessError("coupon:update_status_fail", fmt.Sprintf("更新状态失败"))
		panic(bErr)
	}
	
	// 更新gplutus的account
	//isSuccess := common.NewPlutusService(this.Ctx).DoTransfer(this.user.PlatformId, fmt.Sprintf("cash.coupon.%d", this.RuleId), this.user.GetId(), deductionMoney, bid)
	//if !isSuccess {
	//	eel.Logger.Error(err)
	//	panic(eel.NewBusinessError("coupon:update_account_fail", fmt.Sprintf("使用优惠券更新账户失败")))
	//}

	return nil
}

func (this *Coupon) GetDeductionMoney(orderMoney int) (int, error) {
	conditions := this.GetRule().GetConditions()
	type2condition := make(map[string]*simplejson.Json)
	for _, condition := range conditions{
		type2condition[condition.Type] = condition.Data
	}
	
	// 全额抵扣
	deductionMoney := orderMoney
	// 此券为满减
	//limitReachDeductionInfo, ok := type2condition["reach_deduction"]
	//if ok{
	//	dataStr,_ := json.Marshal(limitReachDeductionInfo)
	//	var aa struct {
	//		Deduction json.Number `json:"deduction"`
	//		Reach json.Number `json:"reach"`
	//	}
	//	_ = json.Unmarshal([]byte(string(dataStr)), &aa)
	//	deductionMoney, _ = aa.Deduction.Float64()
	//	limitMoney, _ := aa.Reach.Float64()
	//	if limitMoney > orderMoney {
	//		eel.Logger.Error("[validation] coupon({}) limited money({}) small than total order money({})", this.Code, limitMoney, orderMoney)
	//		return deductionMoney, eel.NewBusinessError("coupon:order_money_small", fmt.Sprintf("订单金额小于限制的金额"))
	//	}
	//}
	//// 现金抵扣券
	//limitDeductionInfo, ok := type2condition["deduction"]
	//if ok {
	//	limitDeductionInfo = limitDeductionInfo.Get("amount")
	//	deduction := limitDeductionInfo.MustFloat64()
	//	deductionMoney = deduction
	//}
	//// 打折券
	//limitDiscountInfo, ok := type2condition["discount"]
	//if ok {
	//	limitDiscountInfo := limitDiscountInfo.Get("amount")
	//	discount := limitDiscountInfo.MustFloat64()
	//	discountRatio := discount/100
	//	deductionMoney = orderMoney * (1-discountRatio)
	//}
	
	return deductionMoney, nil
}

// 生成优惠券的码
func generateCouponCode(rule *Rule, existCouponCodes []string) string {
	result := snowflakeNode.Generate().String()
	return result
	/*
	generateCode := func(corpId int, ruleId int) string{
		str := "0123456789"
		bytes := []byte(str)
		code := []byte{}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for i := 0; i < 6; i++ {
			code = append(code, bytes[r.Intn(len(bytes))])
		}
		return fmt.Sprintf("%03d%04d%s", corpId, ruleId, string(code))
	}
	couponCode := generateCode(rule.CorpId, rule.Id)
	for _, existCode := range existCouponCodes {
		if existCode == couponCode {
			couponCode = generateCouponCode(rule, existCouponCodes)
		}
	}

	return couponCode
	*/
}

//工厂方法
func NewCoupons(
	ctx context.Context,
	rule *Rule,
	userId int,
	count int,
	sourceType string,
) ([]*Coupon, error) {
	o := eel.GetOrmFromContext(ctx)
	var models []*m_coupon.Coupon
	db := o.Model(&m_coupon.Coupon{}).Where("rule_id", rule.Id).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		bErr := eel.NewBusinessError("coupon:invalid_rule", fmt.Sprintf("错误的优惠券规则"))
		panic(bErr)
		return nil, bErr
	}
	currentCouponCodes := make([]string, 0)
	for _, model := range models {
		currentCouponCodes = append(currentCouponCodes, model.Code)
	}
	dbModels := make([]*m_coupon.Coupon, 0)
	i := 0
	for {
		i++
		if i > count {
			break
		} else {
			couponCode := generateCouponCode(rule, currentCouponCodes)
			currentCouponCodes = append(currentCouponCodes, couponCode)
			dbModels = append(dbModels, &m_coupon.Coupon{
				RuleId: rule.Id,
				UserId: userId,
				Code: couponCode,
				SourceType: m_coupon.COUPON_SOURCE_TYPE_CODE2TYPE[sourceType],
			})
		}
	}
	
	coupons := make([]*Coupon, 0)
	for _, dbModel := range dbModels {
		db := o.Create(dbModel)
		if db.Error != nil {
			eel.Logger.Error(db.Error)
			bErr := eel.NewBusinessError("coupon:create_fail", fmt.Sprintf("生成优惠券失败"))
			panic(bErr)
		}
		
		coupons = append(coupons, NewCouponFromModel(ctx, dbModel))
	}
	//_, err = o.InsertMulti(len(dbModels), dbModels)
	//if err != nil {
	//	eel.Logger.Error(err)
	//	bErr := eel.NewBusinessError("coupon:create_fail", fmt.Sprintf("生成优惠券失败"))
	//	panic(bErr)
	//	return nil, bErr
	//}
	//coupons := make([]*Coupon, 0)
	//for _, dbModel := range dbModels {
	//	coupons = append(coupons, NewCouponFromModel(ctx, dbModel))
	//}
	return coupons, nil
}

//工厂方法
func NewCoupon(
	ctx context.Context,
	userId int,
	sourceType int,
	code string,
	status int,
	orderBid string,
	consumedAt string,
	ruleId int,
) *Coupon {
	o := eel.GetOrmFromContext(ctx)
	model := m_coupon.Coupon{}
	model.UserId = userId
	model.SourceType = sourceType
	model.Code = code
	model.Status = status
	model.OrderBid = orderBid
	model.ConsumedAt = eel.ParseTime(consumedAt)

	model.RuleId = ruleId

	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("coupon:create_fail", fmt.Sprintf("创建失败")))
	}

	return NewCouponFromModel(ctx, &model)
}

//根据model构建对象
func NewCouponFromModel(ctx context.Context, model *m_coupon.Coupon) *Coupon {
	instance := new(Coupon)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.UserId = model.UserId
	instance.SourceType = m_coupon.COUPON_SOURCE_TYPE2CODE[model.SourceType]
	instance.Code = model.Code
	instance.RuleId = model.RuleId
	instance.Status = instance.GetStatus(model)
	instance.OrderBid = model.OrderBid
	instance.ConsumedAt = model.ConsumedAt
	instance.DeductionMoney = model.DeductionMoney
	instance.CreatedAt = model.CreatedAt
	instance.isValidate = false
	instance.isValidityChecked = false

	return instance
}

func init() {
}
