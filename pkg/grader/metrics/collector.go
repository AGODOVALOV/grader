package metrics

import (
	"context"

	"github.com/AGODOVALOV/grader/pkg/client/metrics/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/AGODOVALOV/grader/pkg/client/metrics/metrics"
)

type Collector struct {
	registry *prometheus.Registry
	Metrics  *metrics.CustomMetrics
	Handler  *handler.Handler
}

func NewCollector(ctx context.Context) *Collector {
	reg := prometheus.NewRegistry()

	customMetrics := metrics.NewCustomMetrics()

	reg.MustRegister(collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		//customMetrics.TaskProcessedByStatus,
		//customMetrics.QueueDepthReadyMessages,
		//customMetrics.TaskElapsedTimeBeforeProcess,
		//customMetrics.TaskProcessingDuration,
		//customMetrics.DBOpenConnections,
		//customMetrics.DBIdleConnections,
	)

	return &Collector{
		registry: reg,
		Metrics:  customMetrics,
		Handler:  handler.NewHandler(reg),
	}
}
