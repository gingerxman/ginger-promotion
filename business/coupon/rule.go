package coupon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gingerxman/ginger-promotion/business"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
	"time"
	
	"github.com/gingerxman/eel"
	"github.com/gingerxman/gorm"
)

type Rule struct {
	eel.EntityBase
	Id int
	CorpId int
	Name string
	Count int
	Desc string
	Remark string
	IsActive bool
	IsDeleted bool
	StartDate time.Time
	EndDate time.Time
	CreatedAt time.Time
	Status string
	Conditions []*RuleCondition
	RemainCouponCount int64
	ReceiveStartDate time.Time
	ReceiveEndDate time.Time
	ReceiveStatus string
	//foreign key
}

//Update 更新对象
func (this *Rule) Update(
	name string,
	count int,
	desc string,
	remark string,
	startDate string,
	endDate string,
	receiveStartDate string,
	receiveEndDate string,
	conditions *simplejson.Json,
) error {
	var model m_coupon.Rule
	o := eel.GetOrmFromContext(this.Ctx)

	db:= o.Model(&model).Where("id", this.Id).Update(gorm.Params{
		"name": name,
		"count": count,
		"desc": desc,
		"remark": remark,
		"start_date": eel.ParseTime(startDate),
		"end_date": eel.ParseTime(endDate),
		"receive_start_date": eel.ParseTime(receiveStartDate),
		"receive_end_date": eel.ParseTime(receiveEndDate),
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		bErr := eel.NewBusinessError("rule:update_fail", fmt.Sprintf("更新失败"))
		panic(bErr)
		return bErr
	}
	this.UpdateRuleConditions(conditions)

	return nil
}

//AddRuleCondition 添加限制条件对象
func (this *Rule) AddRuleCondition(ruleCondition *RuleCondition) {
	if this.Conditions == nil {
		this.Conditions = make([]*RuleCondition, 0)
	}

	this.Conditions = append(this.Conditions, ruleCondition)
}

// GetStatus 状态
func (this *Rule,) GetStatus() string{
	now := time.Now()
	startDate := this.StartDate
	endDate := this.EndDate
	var status string
	if (startDate.Before(now) || startDate.Equal(now)) && (endDate.After(now) || endDate.Equal(now)) {
		status = "start"
	} else if endDate.Before(now) {
		status = "end"
	} else {
		status = "not_start"
	}
	return status
}

// GetReceiveStatus 是否在可领的时间范围内
func (this *Rule,) GetReceiveStatus() string{
	now := time.Now()
	startDate := this.ReceiveStartDate
	endDate := this.ReceiveEndDate
	var receiveStatus string
	if (startDate.Before(now) || startDate.Equal(now)) && (endDate.After(now) || endDate.Equal(now)) {
		receiveStatus = "start"
	} else if endDate.Before(now) {
		receiveStatus = "end"
	} else {
		receiveStatus = "not_start"
	}
	return receiveStatus
}

// GetRemainCouponCount 剩余未领的张数
func (this *Rule,) GetRemainCouponCount() int64{
	o := eel.GetOrmFromContext(this.Ctx)
	count, err := o.Model(&m_coupon.Coupon{}).Where("rule_id", this.Id).Count()
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("rule:remain_coupon_count_fail", fmt.Sprintf("获取失败")))
	}
	return int64(this.Count) - count
}

// UpdateRuleConditions 更新规则限制条件
func (this *Rule) UpdateRuleConditions(conditions *simplejson.Json,){
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Where("rule_id", this.Id).Delete(&m_coupon.RuleCondition{})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule_condition:delete_fail", fmt.Sprintf("删除规则限制条件失败")))
	}
	dbModels := make([]*m_coupon.RuleCondition, 0)

	conditionsArr, _ := conditions.Array()
	for _, condition := range conditionsArr{
		cc := condition.(map[string] interface{})
		dataStr,_ := json.Marshal(cc["data"])
		dbModels = append(dbModels, &m_coupon.RuleCondition{
			RuleId: this.Id,
			Type: m_coupon.CONDITION_TYPE_CODE2TYPE[cc["type"].(string)],
			Data: string(dataStr),
		})
	}
	if len(dbModels) != 0{
		// TODO: use o.InsertMulti
		for _, dbModel := range dbModels {
			db := o.Create(dbModel)
			if db.Error != nil {
				eel.Logger.Error(db.Error)
			}
		}
		// _, err = o.InsertMulti(len(dbModels), dbModels)
	}
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule_condition:create_fail", fmt.Sprintf("创建规则限制条件失败")))
	}
}

// AdditionalCount 更新库存
func (this *Rule) AdditionalCount(count int,) {
	o := eel.GetOrmFromContext(this.Ctx)
	var model m_coupon.Rule
	err := o.Model(&m_coupon.Rule{}).Where("id", this.Id).Take(&model)
	if err != nil {
		eel.Logger.Error(err)
		panic(eel.NewBusinessError("rule:get_fail", fmt.Sprintf("获取优惠券规则失败")))
	}
	currentCount := model.Count + count
	db := o.Model(&m_coupon.Rule{}).Where("id", this.Id).Update(gorm.Params{
		"count":currentCount,
	})
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule:update_fail", fmt.Sprintf("更新库存失败")))
	}
}

// Disable 禁用
func (this *Rule) Disable() {
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&m_coupon.Rule{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_active": false,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule:update_fail", fmt.Sprintf("更新规则失败")))
	}
}

// Enable 启用
func (this *Rule) Enable() {
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&m_coupon.Rule{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_active": true,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule:update_fail", fmt.Sprintf("更新规则失败")))
	}
}

// 生成优惠券
func (this *Rule) GenerateCoupons(userId int, count int, sourceType string,) ([]*Coupon, error) {
	return NewCoupons(this.Ctx, this, userId, count, sourceType)
}

func (this *Rule) UpdateRemark (remark string)  {
	o := eel.GetOrmFromContext(this.Ctx)

	db := o.Model(&m_coupon.Rule{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"remark": remark,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule:update_fail", fmt.Sprintf("更新规则失败")))
	}
}

func (this *Rule) GetConditions() []*RuleCondition {
	var models []*m_coupon.RuleCondition
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.RuleCondition{}).Where("rule_id", this.Id).Find(&models)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule:conditions_fail", fmt.Sprintf("获取规则失败")))
	}

	for _, model := range models {
		this.AddRuleCondition(NewRuleConditionFromModel(this.Ctx, model))
	}
	return this.Conditions
}

//工厂方法
func NewRule(
	ctx context.Context,
	corp business.ICorp,
	name string,
	count int,
	desc string,
	remark string,
	startDate string,
	endDate string,
	receiveStartDate string,
	receiveEndDate string,
	conditions *simplejson.Json,
) *Rule {
	o := eel.GetOrmFromContext(ctx)
	model := m_coupon.Rule{}
	model.CorpId = corp.GetId()
	model.Name = name
	model.Count = count
	model.IsActive = true
	model.Desc = desc
	model.Remark = remark
	model.StartDate = eel.ParseTime(startDate)
	model.EndDate = eel.ParseTime(endDate)
	model.ReceiveStartDate = eel.ParseTime(receiveStartDate)
	model.ReceiveEndDate = eel.ParseTime(receiveEndDate)


	db := o.Create(&model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		panic(eel.NewBusinessError("rule:create_fail", fmt.Sprintf("创建优惠券规则失败")))
	}
	rule := NewRuleFromModel(ctx, &model)
	rule.UpdateRuleConditions(conditions)
	// TODO: 创建gplutus的account
	//pUserId := corp.GetRelatedUser().GetId()
	//isSuccess := common.NewPlutusService(ctx).CreateAccount(pUserId, fmt.Sprintf("cash.coupon.%d", rule.Id))
	//if !isSuccess {
	//	eel.Logger.Error(err)
	//	panic(eel.NewBusinessError("rule:create_fail", fmt.Sprintf("优惠券规则账户创建失败")))
	//}
	return rule
}

//根据model构建对象
func NewRuleFromModel(ctx context.Context, model *m_coupon.Rule) *Rule {

	instance := new(Rule)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.Name = model.Name
	instance.Count = model.Count
	instance.Desc = model.Desc
	instance.Remark = model.Remark
	instance.IsActive = model.IsActive
	instance.IsDeleted = model.IsDeleted
	instance.StartDate = model.StartDate
	instance.EndDate = model.EndDate
	instance.CreatedAt = model.CreatedAt
	instance.Status = instance.GetStatus()
	instance.RemainCouponCount = instance.GetRemainCouponCount()
	instance.ReceiveStartDate = model.ReceiveStartDate
	instance.ReceiveEndDate = model.ReceiveEndDate
	instance.ReceiveStatus = instance.GetReceiveStatus()
	instance.Conditions = make([]*RuleCondition, 0)
	return instance
}

func init() {
}
