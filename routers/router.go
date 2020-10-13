package routers

import (
	"github.com/gingerxman/eel"
	"github.com/gingerxman/eel/handler/rest/console"
	"github.com/gingerxman/eel/handler/rest/op"
	"github.com/gingerxman/ginger-promotion/rest/dev"
	"github.com/gingerxman/ginger-promotion/rest/point"
)

func init() {
	eel.RegisterResource(&console.Console{})
	eel.RegisterResource(&op.Health{})
	
	/*
	 point
	 */
	eel.RegisterResource(&point.Product{})
	eel.RegisterResource(&point.DisabledProducts{})
	eel.RegisterResource(&point.Products{})
	eel.RegisterResource(&point.EnabledProducts{})
	
	/*
	 dev
	 */
	eel.RegisterResource(&dev.BDDReset{})
}