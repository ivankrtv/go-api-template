package metrics

import "github.com/prometheus/client_golang/prometheus"

var requestsCountMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
	Help:      "Number of requests recivied",
	Name:      "count",
	Namespace: "service",
	Subsystem: "requests",
}, []string{"method", "service", "code"})

var requestsLatencyMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Help:      "Total duration of requests",
	Name:      "latency",
	Namespace: "service",
	Subsystem: "requests",
}, []string{"method", "service", "code"})

func init() {
	prometheus.MustRegister(requestsCountMetric)
	prometheus.MustRegister(requestsLatencyMetric)
}
