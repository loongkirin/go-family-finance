package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	circuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Current state of the circuit breaker (0: Closed, 1: HalfOpen, 2: Open)",
		},
		[]string{"name"},
	)

	circuitBreakerFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_failures_total",
			Help: "Total number of failures in the circuit breaker",
		},
		[]string{"name"},
	)
)

func init() {
	prometheus.MustRegister(circuitBreakerState)
	prometheus.MustRegister(circuitBreakerFailures)
}
