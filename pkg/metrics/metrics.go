package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)

var (
    KafkaConsumeErrors = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "orders",
        Subsystem: "kafka",
        Name:      "consume_errors_total",
        Help:      "Total number of kafka consume errors",
    })

    DBRetryAttempts = prometheus.NewCounter(prometheus.CounterOpts{
        Namespace: "orders",
        Subsystem: "db",
        Name:      "retry_attempts_total",
        Help:      "Total number of DB retry attempts",
    })

    HandlerRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
        Namespace: "orders",
        Subsystem: "http",
        Name:      "requests_total",
        Help:      "HTTP requests total",
    }, []string{"path", "code"})
)

func init() {
    prometheus.MustRegister(KafkaConsumeErrors, DBRetryAttempts, HandlerRequests)
}


