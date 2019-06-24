package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Monitor struct {
}

var (
	HttpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "audit_ux_http_requests_total",
			Help: "Count of all HTTP requests for portal",
		},
	)
)

func RegisterMonitor() {
	prometheus.MustRegister(HttpRequestsTotal)
}
