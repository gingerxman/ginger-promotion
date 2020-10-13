package point

import (
	"context"
	"github.com/gingerxman/eel"
)

type EncodePointProductService struct {
	eel.ServiceBase
}

func NewEncodePointProductService(ctx context.Context) *EncodePointProductService {
	service := new(EncodePointProductService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodePointProductService) Encode(pointProduct *PointProduct) *RPointProduct {
	if pointProduct == nil {
		return nil
	}
	
	var rProduct *RProduct
	if pointProduct.Product != nil {
		product := pointProduct.Product
		rProduct = &RProduct{
			Id: product.Id,
			SupplierId: product.SupplierId,
			Name: product.Name,
			Thumbnail: product.Thumbnail,
		}
	}

	return &RPointProduct{
		Id: pointProduct.Id,
		CorpId: pointProduct.CorpId,
		PointPrice: pointProduct.PointPrice,
		MoneyPrice: pointProduct.MoneyPrice,
		IsEnabled: pointProduct.IsEnabled,
		BuyLimit: pointProduct.BuyLimit,
		ProductId: pointProduct.ProductId,
		Product: rProduct,
		StartTime: pointProduct.StartTime.Format("2006-01-02 15:05"),
		EndTime: pointProduct.EndTime.Format("2006-01-02 15:05"),
		CreatedAt: pointProduct.CreatedAt.Format("2006-01-02 15:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodePointProductService) EncodeMany(pointProducts []*PointProduct) []*RPointProduct {
	rDatas := make([]*RPointProduct, 0)
	for _, pointProduct := range pointProducts {
		rDatas = append(rDatas, this.Encode(pointProduct))
	}
	
	return rDatas
}

func init() {
}
