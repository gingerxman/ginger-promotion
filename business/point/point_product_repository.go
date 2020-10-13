package point

import (
	"context"
	"github.com/gingerxman/ginger-promotion/business"
	m_point "github.com/gingerxman/ginger-promotion/models/point"
	
	
	"github.com/gingerxman/eel"
)

type PointProductRepository struct {
	eel.RepositoryBase
}

func NewPointProductRepository(ctx context.Context) *PointProductRepository {
	repository := new(PointProductRepository)
	repository.Ctx = ctx
	return repository
}

func (this *PointProductRepository) GetPointProducts(filters eel.Map, orderExprs ...string) []*PointProduct {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_point.PointProduct{})
	
	var models []*m_point.PointProduct
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
	
	instances := make([]*PointProduct, 0)
	for _, model := range models {
		instances = append(instances, NewPointProductFromModel(this.Ctx, model))
	}
	return instances
}

func (this *PointProductRepository) GetPagedPointProducts(filters eel.Map, page *eel.PageInfo, orderExprs ...string) ([]*PointProduct, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	db := o.Model(&m_point.PointProduct{})
	
	var models []*m_point.PointProduct
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
	
	instances := make([]*PointProduct, 0)
	for _, model := range models {
		instances = append(instances, NewPointProductFromModel(this.Ctx, model))
	}
	return instances, paginateResult
}

//GetEnabledPointProductsForCorp 获得启用的PointProduct对象集合
func (this *PointProductRepository) GetEnabledPointProductsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*PointProduct, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	filters["is_enabled"] = true
	
	return this.GetPagedPointProducts(filters, page, "-id")
}


//GetAllPointProductsForCorp 获得所有PointProduct对象集合
func (this *PointProductRepository) GetAllPointProductsForCorp(corp business.ICorp, page *eel.PageInfo, filters eel.Map) ([]*PointProduct, eel.INextPageInfo) {
	filters["corp_id"] = corp.GetId()
	
	return this.GetPagedPointProducts(filters, page, "-id")
}

//GetPointProductInCorp 根据id和corp获得PointProduct对象
func (this *PointProductRepository) GetPointProductInCorp(corp business.ICorp, id int) *PointProduct {
	filters := eel.Map{
		"corp_id": corp.GetId(),
		"id": id,
	}
	
	pointProducts := this.GetPointProducts(filters)
	
	if len(pointProducts) == 0 {
		return nil
	} else {
		return pointProducts[0]
	}
}

//GetPointProduct 根据id和corp获得PointProduct对象
func (this *PointProductRepository) GetPointProduct(id int) *PointProduct {
	filters := eel.Map{
		"id": id,
	}
	
	pointProducts := this.GetPointProducts(filters)
	
	if len(pointProducts) == 0 {
		return nil
	} else {
		return pointProducts[0]
	}
}

func (this *PointProductRepository) GetPointProductByProductIds(ids []int) []*PointProduct {
	filters := eel.Map{
		"product_id__in": ids,
	}
	
	return this.GetPointProducts(filters)
}

func init() {
}
