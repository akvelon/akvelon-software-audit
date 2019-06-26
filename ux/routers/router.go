package routers

import (
	"os"
	"akvelon/akvelon-software-audit/ux/controllers"
	"github.com/astaxie/beego"

	"akvelon/akvelon-software-audit/ux/tracing"
)

func init() {
	var jaeggerSrvName string
	if os.Getenv("JAEGER_SERVICE_NAME") == "" {
		jaeggerSrvName = "akv-ux-jaegger"
	} else {
		jaeggerSrvName = os.Getenv("JAEGER_SERVICE_NAME")
	}

	t, _ := tracing.InitTracer(jaeggerSrvName)

	beego.Router("/", &controllers.MainController{Tracer: t})
	beego.Router("/analyze", &controllers.MainController{Tracer: t}, "post:Analyze")
	beego.Router("/report/:provider/:orgname/:reponame", &controllers.MainController{Tracer: t}, "get:Report")
}
