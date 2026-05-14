package metrics

import "github.com/prometheus/client_golang/prometheus"

type CustomMetrics struct {
	TaskProcessedByStatus *prometheus.CounterVec
	HTTPRequestTotal      *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	UploadTotal           *prometheus.CounterVec
	UploadFileSize        *prometheus.HistogramVec
	ReviewCreatedTotal    *prometheus.CounterVec
	LoginAttemptsTotal    *prometheus.CounterVec
	RegisterAttemptsTotal *prometheus.CounterVec
}

func NewCustomMetrics() *CustomMetrics {
	return &CustomMetrics{
		TaskProcessedByStatus: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_tasks_processed_total",
				Help: "Number of tasks processed per status",
			},
			[]string{"status"},
		),

		HTTPRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),

		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "client_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),

		UploadTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_uploads_total",
				Help: "Total number of file upload attempts",
			},
			[]string{"result"},
		),

		UploadFileSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "client_upload_file_size_bytes",
				Help: "Uploaded file size in bytes",
				Buckets: []float64{
					1024,
					10 * 1024,
					100 * 1024,
					500 * 1024,
					1024 * 1024,
					5 * 1024 * 1024,
				},
			},
			[]string{"result"},
		),

		ReviewCreatedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_reviews_created_total",
				Help: "Total number of review creation attempts",
			},
			[]string{"result"},
		),

		LoginAttemptsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_login_attempts_total",
				Help: "Total number of login attempts",
			},
			[]string{"result"},
		),

		RegisterAttemptsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "client_register_attempts_total",
				Help: "Total number of register attempts",
			},
			[]string{"result"},
		),
	}
}
