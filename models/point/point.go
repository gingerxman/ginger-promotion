package point

import (
	"github.com/gingerxman/eel"
	"time"
)

//Product Model
type PointProduct struct {
	eel.Model
	CorpId int `gorm:"index"` //foreign key for corp
	ProductId int `gorm:"index"`
	PointPrice int // 积分价格
	MoneyPrice int // 现金价格
	StartTime time.Time
	EndTime time.Time
	BuyLimit int // 兑换限制
	IsEnabled bool `gorm:"default:true"`
}
func (self *PointProduct) TableName() string {
	return "point_product"
}



func init() {
	eel.RegisterModel(new(PointProduct))
}
