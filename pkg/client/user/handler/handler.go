package handler

import (
	"html/template"

	"github.com/AGODOVALOV/grader/pkg/client/metrics"
	"github.com/AGODOVALOV/grader/pkg/client/user/usecase"
)

// UserHandler handles user-related operations.
type UserHandler struct {
	template         *template.Template
	Service          *usecase.UserService
	MetricsCollector *metrics.Collector
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(t *template.Template, s *usecase.UserService, metricsCollector *metrics.Collector) *UserHandler {
	return &UserHandler{
		template:         t,
		Service:          s,
		MetricsCollector: metricsCollector,
	}
}
