package point

import (
	"context"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-promotion/business/product"
)

type FillPointProductService struct {
	eel.ServiceBase
}

func NewFillPointProductService(ctx context.Context) *FillPointProductService {
	service := new(FillPointProductService)
	service.Ctx = ctx
	return service
}

func (this *FillPointProductService) Fill(pointProducts []*PointProduct, option eel.FillOption) {
	if len(pointProducts) == 0 {
		return
	}
	
	ids := make([]int, 0)
	for _, product := range pointProducts {
		ids = append(ids, product.Id)
	}

	this.fillProduct(pointProducts, ids)
	return
}


func (this *FillPointProductService) fillProduct(pointProducts []*PointProduct, ids []int) {
	//获取product id集合
	productIds := make([]int, 0)
	for _, pointProduct := range pointProducts {
		productIds = append(productIds, pointProduct.ProductId)
	}
	
	// 获取products
	products := product.NewProductRepository(this.Ctx).GetProducts(productIds)
	
	// 构建<id, product>
	id2product := make(map[int]*product.Product)
	for _, product := range products {
		id2product[product.Id] = product
	}
	
	for _, pointProduct := range pointProducts {
		if product, ok := id2product[pointProduct.ProductId]; ok {
			pointProduct.Product = product
		}
	}
}


func init() {
}
