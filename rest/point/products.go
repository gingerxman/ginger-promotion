package point

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-promotion/business/account"
	b_point "github.com/gingerxman/ginger-promotion/business/point"
	//b_point "github.com/gingerxman/ginger-promotion/models/point"
)

type Products struct {
	eel.RestResource
}

func (this *Products) Resource() string {
	return "point.products"
}

func (this *Products) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?product_ids:json-array", "?filters:json", "?fill_options:json-array"},
	}
}

func (this *Products) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	pageInfo := req.GetPageInfo()
	filters := req.GetOrmFilters()
	
	productIds := req.GetIntArray("product_ids")
	var pointProducts []*b_point.PointProduct
	var nextPageInfo eel.INextPageInfo
	if len(productIds) == 0 {
		// 获取ids指定的积分商品集合
		corp := account.GetCorpFromContext(bCtx)
		pointProducts, nextPageInfo = b_point.NewPointProductRepository(bCtx).GetAllPointProductsForCorp(corp, pageInfo, filters)
	} else {
		// 获取分页商品集合
		pointProducts = b_point.NewPointProductRepository(bCtx).GetPointProductByProductIds(productIds)
		nextPageInfo = eel.MockPaginate(0, pageInfo)
	}
	
	b_point.NewFillPointProductService(bCtx).Fill(pointProducts, eel.FillOption{})
	rows := b_point.NewEncodePointProductService(bCtx).EncodeMany(pointProducts)
	
	ctx.Response.JSON(eel.Map{
		"products": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}

