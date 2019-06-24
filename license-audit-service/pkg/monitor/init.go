package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Monitor struct {
}

var (
	httpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "audit_srv_http_requests_total",
			Help: "Count of all HTTP requests for audit service",
		})
)

func (m *Monitor) RegisterMonitor() {
	prometheus.MustRegister(httpRequestsTotal)
}

func (m *Monitor) GetHttpRequestsTotal() prometheus.Counter {
	return httpRequestsTotal
}
