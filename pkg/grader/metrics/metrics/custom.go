package metrics

import "github.com/prometheus/client_golang/prometheus"

type CustomMetrics struct {
	HTTPRequestTotal    *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec

	JobsReceivedTotal  *prometheus.CounterVec
	JobsProcessedTotal *prometheus.CounterVec
	JobDuration        *prometheus.HistogramVec

	DockerRunsTotal   *prometheus.CounterVec
	DockerRunDuration *prometheus.HistogramVec

	S3DownloadsTotal *prometheus.CounterVec

	CallbacksTotal *prometheus.CounterVec

	WorkerQueueDepth prometheus.Gauge
}

func NewCustomMetrics() *CustomMetrics {
	return &CustomMetrics{
		HTTPRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grader_http_requests_total",
				Help: "Total number of grader HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grader_http_request_duration_seconds",
				Help:    "Grader HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),

		JobsReceivedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grader_jobs_received_total",
				Help: "Total number of grader jobs received",
			},
			[]string{"task"},
		),

		JobsProcessedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grader_jobs_processed_total",
				Help: "Total number of grader jobs processed",
			},
			[]string{"task", "result"},
		),

		JobDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "grader_job_duration_seconds",
				Help: "Total grader job processing duration in seconds",
				Buckets: []float64{
					0.1,
					0.5,
					1,
					2,
					5,
					10,
					30,
					60,
					120,
				},
			},
			[]string{"task", "result"},
		),

		DockerRunsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grader_docker_runs_total",
				Help: "Total number of docker test runs",
			},
			[]string{"task"},
		),

		DockerRunDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "grader_docker_run_duration_seconds",
				Help: "Docker test run duration in seconds",
				Buckets: []float64{
					0.1,
					0.5,
					1,
					2,
					5,
					10,
					30,
					60,
				},
			},
			[]string{"task"},
		),

		S3DownloadsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grader_s3_downloads_total",
				Help: "Total number of S3 download attempts",
			},
			[]string{"task", "result"},
		),

		CallbacksTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grader_callbacks_total",
				Help: "Total number of callback requests",
			},
			[]string{"task", "result"},
		),

		WorkerQueueDepth: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "grader_worker_queue_depth",
				Help: "Current number of jobs waiting in grader worker queue",
			},
		),
	}
}
