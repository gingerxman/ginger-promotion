package point

type RProduct struct {
	Id int `json:"id"`
	SupplierId int `json:"supplier_id"`
	Name string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type RPointProduct struct {
	Id int `json:"id"`
	CorpId int `json:"corp_id"`
	ProductId int `json:"product_id"`
	PointPrice int `json:"point_price"`
	MoneyPrice int `json:"money_price"`
	IsEnabled bool `json:"is_enabled"`
	BuyLimit int `json:"buy_limit"`
	Product *RProduct `json:"product"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	CreatedAt string `json:"created_at"`
}

func init() {
}
