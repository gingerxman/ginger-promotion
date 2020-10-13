package point

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-promotion/business/account"
	"github.com/gingerxman/ginger-promotion/business/point"
)

type EnabledProducts struct {
	eel.RestResource
}

func (this *EnabledProducts) Resource() string {
	return "point.enabled_products"
}

func (this *EnabledProducts) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"?filters:json", "?fill_options:json-array"},
	}
}

func (this *EnabledProducts) Get(ctx *eel.Context) {
	req := ctx.Request
	bCtx := ctx.GetBusinessContext()
	pageInfo := req.GetPageInfo()
	filters := req.GetOrmFilters()
	
	corp := account.GetCorpFromContext(bCtx)
	pointProducts, nextPageInfo := point.NewPointProductRepository(bCtx).GetEnabledPointProductsForCorp(corp, pageInfo, filters)
	
	point.NewFillPointProductService(bCtx).Fill(pointProducts, eel.FillOption{})
	rows := point.NewEncodePointProductService(bCtx).EncodeMany(pointProducts)
	
	ctx.Response.JSON(eel.Map{
		"products": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}

