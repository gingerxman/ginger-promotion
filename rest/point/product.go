package point

import (
	"fmt"
	"github.com/gingerxman/eel"
	"github.com/gingerxman/ginger-promotion/business/account"
	"github.com/gingerxman/ginger-promotion/business/point"
)

type Product struct {
	eel.RestResource
}

func (this *Product) Resource() string {
	return "point.product"
}

func (this *Product) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"bid", "?with_options:json"},
		"PUT": []string{
			"product_id:int",
			"point_price:int",
			"money_price:int",
			"buy_limit:int",
			"start_time",
			"end_time",
		},
		"POST": []string{
			"id:int",
			"point_price:int",
			"money_price:int",
			"buy_limit:int",
			"start_time",
			"end_time",
		},
		"DELETE": []string{
			"id:int",
		},
	}
}

func (this *Product) Get(ctx *eel.Context) {
	//get order
	// req := ctx.Request
	// bCtx := ctx.GetBusinessContext()
	// bid := req.GetString("bid")
	//order := b_order.NewProductRepository(bCtx).GetProductByBid(bid)
	//
	////处理fill options
	//fillOptions := req.GetJSON("with_options")
	//fillService := b_order.NewFillProductService(bCtx)
	//if fillOptions == nil{
	//	fillOptions = eel.Map{}
	//}
	////处理invoice的fill options
	//invoiceFillOptions := eel.Map{
	//	"with_products": true,
	//}
	////if option, ok := fillOptions["with_settlements"]; ok {
	////	invoiceFillOptions["with_settlements"] = option
	////} else {
	////	invoiceFillOptions["with_settlements"] = true
	////}
	//fillOptions["with_invoice"] = invoiceFillOptions
	//fillOptions["with_operation_log"] = true
	//fillService.Fill([]*b_order.Product{order}, fillOptions)
	//
	////encode
	//data := b_order.NewEncodeProductService(bCtx).Encode(order)
	
	ctx.Response.JSON(eel.Map{})
}

func (this *Product) Put(ctx *eel.Context) {
	req := ctx.Request
	productId, _ := req.GetInt("product_id")
	pointPrice, _ := req.GetInt("point_price")
	moneyPrice, _ := req.GetInt("money_price")
	buyLimit, _ := req.GetInt("buy_limit")
	startTime := eel.ParseTime(req.GetString("start_time"))
	endTime := eel.ParseTime(req.GetString("end_time"))
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	pointProduct := point.NewPointProduct(bCtx, corp, productId, pointPrice, moneyPrice, buyLimit, startTime, endTime)
	
	ctx.Response.JSON(eel.Map{
		"id": pointProduct.Id,
	})
}

func (this *Product) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	pointPrice, _ := req.GetInt("point_price")
	moneyPrice, _ := req.GetInt("money_price")
	buyLimit, _ := req.GetInt("buy_limit")
	startTime := eel.ParseTime(req.GetString("start_time"))
	endTime := eel.ParseTime(req.GetString("end_time"))
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	pointProduct := point.NewPointProductRepository(bCtx).GetPointProductInCorp(corp, id)
	if pointProduct == nil {
		ctx.Response.Error("product:invalid_point_product", fmt.Sprintf("id(%d)", id))
		return
	}
	
	if pointProduct.IsEnabled {
		ctx.Response.Error("product:update_enabled_point_product", fmt.Sprintf("id(%d)", id))
		return
	}
	
	pointProduct.Update(pointPrice, moneyPrice, buyLimit, startTime, endTime)
	
	ctx.Response.JSON(eel.Map{
	})
}

func (this *Product) Delete(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	
	bCtx := ctx.GetBusinessContext()
	corp := account.GetCorpFromContext(bCtx)
	pointProduct := point.NewPointProductRepository(bCtx).GetPointProductInCorp(corp, id)
	if pointProduct == nil {
		ctx.Response.Error("product:invalid_point_product", fmt.Sprintf("id(%d)", id))
		return
	}
	
	if pointProduct.IsEnabled {
		ctx.Response.Error("product:delete_enabled_point_product", fmt.Sprintf("id(%d)", id))
		return
	}
	
	err := pointProduct.Delete()
	if err != nil {
		eel.Logger.Error(err)
	}
	
	ctx.Response.JSON(eel.Map{
	})
}
