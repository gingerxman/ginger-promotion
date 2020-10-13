package point

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-promotion/business"
	m_point "github.com/gingerxman/ginger-promotion/models/point"
	"github.com/gingerxman/ginger-promotion/business/product"
	"github.com/gingerxman/gorm"
	"time"
)

type PointProduct struct {
	eel.EntityBase
	Id int
	CorpId int
	
	PointPrice int
	MoneyPrice int
	IsEnabled bool
	BuyLimit int
	
	ProductId int
	Product *product.Product
	
	StartTime time.Time
	EndTime time.Time
	
	CreatedAt time.Time
}

func (this *PointProduct) enable(isEnable bool) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	db := o.Model(&m_point.PointProduct{}).Where(eel.Map{
		"id": this.Id,
		"corp_id": this.CorpId,
	}).Update(gorm.Params{
		"is_enabled": isEnable,
	})

	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *PointProduct) Enable() {
	this.enable(true)
}

func (this *PointProduct) Disable() {
	this.enable(false)
}

func (this *PointProduct) Update(
	pointPrice int,
	moneyPrice int,
	buyLimit int,
	startTime time.Time,
	endTime time.Time,
) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_point.PointProduct{}).Where("id", this.Id).Update(gorm.Params{
		"point_price": pointPrice,
		"money_price": moneyPrice,
		"buy_limit": buyLimit,
		"start_time": startTime,
		"end_time": endTime,
	})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
	}
}

func (this *PointProduct) Delete() error {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Where("id", this.Id).Delete(&m_point.PointProduct{})
	
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return db.Error
	}
	
	return nil
}

func NewPointProduct(
	ctx context.Context,
	corp business.ICorp,
	productId int,
	pointPrice int,
	moneyPrice int,
	buyLimit int,
	startTime time.Time,
	endTime time.Time,
) *PointProduct {
	o := eel.GetOrmFromContext(ctx)
	model := &m_point.PointProduct{
		CorpId: corp.GetId(),
		ProductId: productId,
		PointPrice: pointPrice,
		MoneyPrice: moneyPrice,
		BuyLimit: buyLimit,
		StartTime: startTime,
		EndTime: endTime,
	}
	
	db := o.Create(model)
	if db.Error != nil {
		eel.Logger.Error(db.Error)
		return nil
	}
	
	return NewPointProductFromModel(ctx, model)
}

//根据model构建对象
func NewPointProductFromModel(ctx context.Context, model *m_point.PointProduct) *PointProduct {
	instance := new(PointProduct)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.CorpId = model.CorpId
	instance.ProductId = model.ProductId
	instance.PointPrice = model.PointPrice
	instance.MoneyPrice = model.MoneyPrice
	instance.IsEnabled = model.IsEnabled
	instance.BuyLimit = model.BuyLimit
	
	instance.StartTime = model.StartTime
	instance.EndTime = model.EndTime
	instance.CreatedAt = model.CreatedAt

	return instance
}

func init() {
}
