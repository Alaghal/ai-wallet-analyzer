package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HTTPRequestTotal    *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	HTTPInFlight        prometheus.Gauge

	ProviderRequestsTotal   *prometheus.CounterVec
	ProviderRequestDuration *prometheus.HistogramVec

	LLMRequestsTotal   *prometheus.CounterVec
	LLMRequestDuration *prometheus.HistogramVec
}

func MustNew() *Metrics {
	m := &Metrics{
		HTTPRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "wallet_analyzer",
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests processed by the service.",
			},
			[]string{"method", "route", "status"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "wallet_analyzer",
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "HTTP request duration in seconds.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "route", "status"},
		),
		HTTPInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "wallet_analyzer",
				Subsystem: "http",
				Name:      "in_flight_requests",
				Help:      "Current number of in-flight HTTP requests.",
			},
		),
		ProviderRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "wallet_analyzer",
				Subsystem: "provider",
				Name:      "requests_total",
				Help:      "Total number of wallet activity provider requests.",
			},
			[]string{"provider", "status"},
		),
		ProviderRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "wallet_analyzer",
				Subsystem: "provider",
				Name:      "request_duration_seconds",
				Help:      "Wallet activity provider request duration in seconds.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"provider", "status"},
		),
		LLMRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "wallet_analyzer",
				Subsystem: "llm",
				Name:      "requests_total",
				Help:      "Total number of LLM requests.",
			},
			[]string{"provider", "status"},
		),
		LLMRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "wallet_analyzer",
				Subsystem: "llm",
				Name:      "request_duration_seconds",
				Help:      "LLM request duration in seconds.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"provider", "status"},
		),
	}

	prometheus.MustRegister(
		m.HTTPRequestTotal,
		m.HTTPRequestDuration,
		m.HTTPInFlight,
		m.ProviderRequestsTotal,
		m.ProviderRequestDuration,
		m.LLMRequestsTotal,
		m.LLMRequestDuration,
	)

	return m
}
