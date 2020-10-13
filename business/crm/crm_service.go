package crm

import (
	"context"
	"github.com/gingerxman/eel"
)

type CrmService struct {
	eel.RepositoryBase
}

func NewCrmService(ctx context.Context) *CrmService {
	repository := new(CrmService)
	repository.Ctx = ctx
	return repository
}

func (this *CrmService) RecordOrder(bid string, money int) {
	_, err := eel.NewResource(this.Ctx).Put("ginger-crm", "point.finished_order", eel.Map{
		"bid": bid,
		"money": money,
	})
	
	if err != nil {
		eel.Logger.Error(err)
	}
}

func init() {
}
