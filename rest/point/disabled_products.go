package point

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-promotion/business/account"
	"github.com/gingerxman/ginger-promotion/business/point"
)

type DisabledProducts struct {
	eel.RestResource
}

func (this *DisabledProducts) Resource() string {
	return "point.disabled_products"
}

func (this *DisabledProducts) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{"ids:json-array"},
		"DELETE": []string{"ids:json-array"},
	}
}

func (this *DisabledProducts) Put(ctx *eel.Context) {
	req := ctx.Request
	ids := req.GetIntArray("ids")
	
	for _, id := range ids {
		bCtx := ctx.GetBusinessContext()
		corp := account.GetCorpFromContext(bCtx)
		pointProduct := point.NewPointProductRepository(bCtx).GetPointProductInCorp(corp, id)
		
		if pointProduct == nil {
			// ctx.Response.Error("disabled_product:invalid_point_product", fmt.Sprintf("id=%d", id))
			// return
		} else {
			pointProduct.Disable()
		}
	}
	
	ctx.Response.JSON(eel.Map{})
}

func (this *DisabledProducts) Delete(ctx *eel.Context) {
	req := ctx.Request
	ids := req.GetIntArray("ids")
	
	for _, id := range ids {
		bCtx := ctx.GetBusinessContext()
		corp := account.GetCorpFromContext(bCtx)
		pointProduct := point.NewPointProductRepository(bCtx).GetPointProductInCorp(corp, id)
		
		if pointProduct == nil {
			// ctx.Response.Error("disabled_product:invalid_point_product", fmt.Sprintf("id=%d", id))
			// return
		} else {
			pointProduct.Enable()
		}
	}
	
	ctx.Response.JSON(eel.Map{})
}
