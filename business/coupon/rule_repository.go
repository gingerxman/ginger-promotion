package coupon

import (
	"context"
	"github.com/gingerxman/ginger-promotion/business"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"

	
	"github.com/gingerxman/eel"
)

type RuleRepository struct {
	eel.RepositoryBase
}

func NewRuleRepository(ctx context.Context) *RuleRepository {
	repository := new(RuleRepository)
	repository.Ctx = ctx
	return repository
}

func (this *RuleRepository) GetRules(filters eel.Map, orderExprs ...string) []*Rule {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.Rule{})
	
	var models []*m_coupon.Rule
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
	
	Rules := make([]*Rule, 0)
	for _, model := range models {
		Rules = append(Rules, NewRuleFromModel(this.Ctx, model))
	}
	return Rules
}

func (this *RuleRepository) GetPagedRules(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*Rule, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_coupon.Rule{})
	
	var models []*m_coupon.Rule
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
	
	Rules := make([]*Rule, 0)
	for _, model := range models {
		Rules = append(Rules, NewRuleFromModel(this.Ctx, model))
	}
	return Rules, paginateResult
}

//GetRulesForCorp 获得启用的Rule对象集合
func (this *RuleRepository) GetRulesForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*Rule, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedRules(filters, page, "-id")
}

//GetEnabledRules 获得启用的Rule对象集合
func (this *RuleRepository) GetEnabledRules(page *eel.PageInfo, filters eel.Map) ([]*Rule, eel.INextPageInfo) {
	//filters["corp_id"] = corp.GetId()
	filters["is_deleted"] = false

	return this.GetPagedRules(filters, page, "-id")
}

//GetAllRulesForCorp 获得所有Rule对象集合
func (this *RuleRepository) GetAllRulesForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*Rule, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedRules(filters, page, "-id")
}

//GetRuleInCorp 根据id和corp获得Rule对象
func (this *RuleRepository) GetRuleInCorp(corp business.ICorp, id int) *Rule {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	Rules := this.GetRules(filters)
	
	if len(Rules) == 0 {
		return nil
	} else {
		return Rules[0]
	}
}

//GetRule 根据id获得Rule对象
func (this *RuleRepository) GetRule(id int) *Rule {
	filters := eel.Map{
		"id": id,
	}
	
	Rules := this.GetRules(filters)
	
	if len(Rules) == 0 {
		return nil
	} else {
		return Rules[0]
	}
}

//GetRuleByName 根据name获得Rule对象
func (this *RuleRepository) GetRuleByName(name string) *Rule {
	filters := eel.Map{
		"name": name,
	}
	
	Rules := this.GetRules(filters)
	
	if len(Rules) == 0 {
		return nil
	} else {
		return Rules[0]
	}
}

func init() {
}
