package routers

import (
	"akvelon/akvelon-software-audit/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/analyze", &controllers.MainController{}, "post:Analyze")
}
