package params

import (
	"encoding/json"
	
	"github.com/gingerxman/eel"
)

type CouponInfo struct {
	Code string
	Money float64
}
type OrderInfo struct {
	Bid string
	Money float64
	UserId int64 `json:"user_id"`
}
type ProductInfo struct {
	Id int
}

type CouponParams struct {
	Coupon *CouponInfo
	Order *OrderInfo
	Products []*ProductInfo
}

func NewParseCouponParams(couponInfoStr string, orderInfoStr string, products []interface{}) *CouponParams{
	couponParams := new(CouponParams)
	couponInfo := CouponInfo{}
	if couponInfoStr != "" {
		err := json.Unmarshal([]byte(couponInfoStr), &couponInfo)
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("coupon:parse_coupon_info_fail", "解析CouponInfo出错"))
		}
		couponParams.Coupon = &couponInfo
	}
	orderInfo := OrderInfo{}
	if orderInfoStr != "" {
		err := json.Unmarshal([]byte(orderInfoStr), &orderInfo)
		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("coupon:parse_order_info_fail", "解析OrderInfo出错"))
		}
		couponParams.Order = &orderInfo
	}

	productInfos := make([]*ProductInfo, 0)
	for _, product := range products {
		data := product.(map[string]interface{})
		id, _ := data["id"].(json.Number).Int64()
		productInfo := ProductInfo{}
		productInfo.Id = int(id)
		productInfos = append(productInfos, &productInfo)
	}
	couponParams.Products = productInfos

	return couponParams
}