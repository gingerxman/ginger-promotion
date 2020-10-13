package coupon

import (
	"context"
	"github.com/bitly/go-simplejson"
	m_coupon "github.com/gingerxman/ginger-promotion/models/coupon"
	"time"

	"github.com/gingerxman/eel"
)

type RuleCondition struct {
	eel.EntityBase
	Id int
	Type string
	Data *simplejson.Json
	CreatedAt time.Time

	//foreign key
	RuleId int //refer to rule
}

//根据model构建对象
func NewRuleConditionFromModel(ctx context.Context, model *m_coupon.RuleCondition) *RuleCondition {
	js, _ := simplejson.NewJson([]byte(model.Data))
	instance := new(RuleCondition)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.Type = m_coupon.CONDITION_TYPE2CODE[model.Type]
	instance.Data = js
	instance.RuleId = model.RuleId
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
