package metrics

import (
	"context"

	"github.com/AGODOVALOV/grader/pkg/grader/metrics/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/AGODOVALOV/grader/pkg/grader/metrics/metrics"
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
		customMetrics.HTTPRequestTotal,
		customMetrics.HTTPRequestDuration,
		customMetrics.JobsReceivedTotal,
		customMetrics.JobsProcessedTotal,
		customMetrics.JobDuration,
		customMetrics.DockerRunsTotal,
		customMetrics.DockerRunDuration,
		customMetrics.S3DownloadsTotal,
		customMetrics.CallbacksTotal,
		customMetrics.WorkerQueueDepth,
	)

	return &Collector{
		registry: reg,
		Metrics:  customMetrics,
		Handler:  handler.NewHandler(reg),
	}
}
