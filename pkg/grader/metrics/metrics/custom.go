package metrics

import "github.com/prometheus/client_golang/prometheus"

type CustomMetrics struct {
	TaskProcessedByStatus        *prometheus.CounterVec
	QueueDepthReadyMessages      *prometheus.GaugeVec
	TaskElapsedTimeBeforeProcess *prometheus.HistogramVec
	TaskProcessingDuration       *prometheus.HistogramVec
	DBOpenConnections            prometheus.Gauge
	DBIdleConnections            prometheus.Gauge
}

func NewCustomMetrics() *CustomMetrics {
	return &CustomMetrics{
		TaskProcessedByStatus: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "qproc_tasks_processed_total",
				Help: "Number of tasks processed per status",
			},
			[]string{"status"},
		),
		//QueueDepthReadyMessages: prometheus.NewGaugeVec(
		//	prometheus.GaugeOpts{
		//		Name: "qproc_queue_messages_ready",
		//		Help: "Number of ready messages in the queue",
		//	},
		//	[]string{"queue"},
		//),
		//TaskElapsedTimeBeforeProcess: prometheus.NewHistogramVec(
		//	prometheus.HistogramOpts{
		//		Name: "qproc_task_elapsed_time_before_process",
		//		Help: "Time elapsed before task processing",
		//	},
		//	[]string{"queue"},
		//),
		//TaskProcessingDuration: prometheus.NewHistogramVec(
		//	prometheus.HistogramOpts{
		//		Name: "qproc_task_processing_duration_seconds",
		//		Help: "Time processing task",
		//	},
		//	[]string{"code", "queue"},
		//),
		//DBOpenConnections: prometheus.NewGauge(
		//	prometheus.GaugeOpts{
		//		Name: "qproc_db_open_connections",
		//		Help: "Number of open database connections",
		//	},
		//),
		//DBIdleConnections: prometheus.NewGauge(
		//	prometheus.GaugeOpts{
		//		Name: "qproc_db_idle_connections",
		//		Help: "Number of idle database connections",
		//	},
		//),
	}
}
