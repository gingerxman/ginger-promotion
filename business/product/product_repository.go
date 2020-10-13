package product

import (
	"context"
	"encoding/json"
	"github.com/gingerxman/eel"
)

type ProductRepository struct {
	eel.RepositoryBase
}

func NewProductRepository(ctx context.Context) *ProductRepository {
	repository := new(ProductRepository)
	repository.Ctx = ctx
	return repository
}
func (this *ProductRepository) makeProducts(productDatas []interface{}) []*Product {
	products := make([]*Product, 0)
	for _, d := range productDatas {
		productData := d.(map[string]interface{})
		
		isDeleted := productData["is_deleted"].(bool)
		productId64, _ := productData["id"].(json.Number).Int64()
		corpId64, _ := productData["corp_id"].(json.Number).Int64()
		
		baseInfo := productData["base_info"].(map[string]interface{})
		rawProductId64, _ := baseInfo["id"].(json.Number).Int64()
		supplierId64, _ := baseInfo["supplier_id"].(json.Number).Int64()
		
		logisticsInfo := productData["logistics_info"].(map[string]interface{})
		unifiedPostageMoney64, _ := logisticsInfo["unified_postage_money"].(json.Number).Int64()
		
		product := &Product{
			Id: int(productId64),
			RawProductId: int(rawProductId64),
			CorpId: int(corpId64),
			SupplierId: int(supplierId64),
			Status: baseInfo["shelf_status"].(string),
			Name: baseInfo["name"].(string),
			Thumbnail: baseInfo["thumbnail"].(string),
			
			IsDeleted: isDeleted,
			LogisticsInfo: &sLogisticsInfo{
				PostageType: logisticsInfo["postage_type"].(string),
				UnifiedPostageMoney: int(unifiedPostageMoney64),
				LimitZoneType: logisticsInfo["limit_zone_type_code"].(string),
			},
		}
		
		skuDatas := productData["skus"].([]interface{})
		for _, d := range skuDatas {
			skuData := d.(map[string]interface{})
			id64, _ := skuData["id"].(json.Number).Int64()
			price64, _ := skuData["price"].(json.Number).Int64()
			stocks64, _ := skuData["stocks"].(json.Number).Int64()
			product.Skus = append(product.Skus, &sku{
				Id: int(id64),
				Name: skuData["name"].(string),
				DisplayName: skuData["display_name"].(string),
				Code: skuData["code"].(string),
				Price: int(price64),
				Stocks: int(stocks64),
			})
			
		}
		products = append(products, product)
	}
	
	return products
}

func (this *ProductRepository) GetProducts(ids []int) []*Product {
	options := []string{"with_sku", "with_logistics"}
	resp, err := eel.NewResource(this.Ctx).Get("ginger-product", "product.products", eel.Map{
		"ids": eel.ToJsonString(ids),
		"fill_options": eel.ToJsonString(options),
	})
	
	if err != nil {
		eel.Logger.Error(err)
		return nil
	}
	
	respData := resp.Data()
	productDatas := respData.Get("products")
	return this.makeProducts(productDatas.MustArray())
}

func (this *ProductRepository) UseSkuStocks(skuId int, count int) error {
	resp, err := eel.NewResource(this.Ctx).Put("ginger-product", "product.sku_stock_consumption", eel.Map{
		"sku_id": skuId,
		"count": count,
	})
	
	if err != nil {
		eel.Logger.Error(resp)
		eel.Logger.Error(err)
		return err
	}
	return nil
}

func (this *ProductRepository) AddSkuStocks(skuId int, count int) error {
	resp, err := eel.NewResource(this.Ctx).Delete("ginger-product", "product.sku_stock_consumption", eel.Map{
		"sku_id": skuId,
		"count": count,
	})
	
	if err != nil {
		eel.Logger.Error(resp)
		eel.Logger.Error(err)
		return err
	}
	return nil
}

func init() {
}
