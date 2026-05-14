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
		customMetrics.TaskProcessedByStatus,
		customMetrics.HTTPRequestTotal,
		customMetrics.HTTPRequestDuration,
		customMetrics.UploadTotal,
		customMetrics.UploadFileSize,
		customMetrics.ReviewCreatedTotal,
		customMetrics.LoginAttemptsTotal,
		customMetrics.RegisterAttemptsTotal,
	)

	return &Collector{
		registry: reg,
		Metrics:  customMetrics,
		Handler:  handler.NewHandler(reg),
	}
}
