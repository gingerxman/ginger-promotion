package dev

import (
	"github.com/gingerxman/eel"
)

type BDDReset struct {
	eel.RestResource
}

func (this *BDDReset) Resource() string {
	return "dev.bdd_reset"
}

func (this *BDDReset) SkipAuthCheck() bool {
	return true
}

func (r *BDDReset) IsForDevTest() bool {
	return true
}

func (this *BDDReset) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT":  []string{},
	}
}

func (this *BDDReset) Put(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	o := eel.GetOrmFromContext(bCtx)
	
	o.Exec("delete from point_product")
	
	ctx.Response.JSON(eel.Map{})
}

