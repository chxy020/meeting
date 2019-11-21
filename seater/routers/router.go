package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"

	"seater/controllers"
)

func init() {
	ns := beego.NewNamespace("v1",
		beego.NSOptions("/*", func(ctx *context.Context) {
			ctx.Output.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT, PATCH")
			ctx.Output.Header("Access-Control-Allow-Headers", "accept, content-type")
			ctx.Output.Body([]byte("."))
		}),
		beego.NSNamespace("/meetings",
			beego.NSInclude(&controllers.MeetingController{}),
		),
	)
	beego.AddNamespace(ns)
}
