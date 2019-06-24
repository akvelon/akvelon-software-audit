package main

import (
	_ "akvelon/akvelon-software-audit/ux/routers"
	"akvelon/akvelon-software-audit/ux/monitor"

	"github.com/astaxie/beego"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	beego.Handler("/metrics", promhttp.Handler())
	monitor.RegisterMonitor()

	beego.Run()
}
