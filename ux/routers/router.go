package routers

import (
	"akvelon/akvelon-software-audit/ux/controllers"
	"github.com/astaxie/beego"

	"akvelon/akvelon-software-audit/ux/tracing"
)

func init() {
	t, _ := tracing.InitTracer(beego.AppConfig.String("jaeger-srv-name"))

	beego.Router("/", &controllers.MainController{Tracer: t})
	beego.Router("/analyze", &controllers.MainController{Tracer: t}, "post:Analyze")
	beego.Router("/report/:provider/:orgname/:reponame", &controllers.MainController{Tracer: t}, "get:Report")
}
