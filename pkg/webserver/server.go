// Package webserver provides an HTTP server.
package webserver

import (
	"context"
	"embed"
	"html/template"
	"net"
	"net/http"
	"strconv"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/webserver/config"
	"github.com/AGODOVALOV/grader/pkg/webserver/handler"
	"github.com/AGODOVALOV/grader/pkg/webserver/middleware"
	"github.com/swaggo/http-swagger"
)

//go:embed templates/*.html
var templateFS embed.FS

// Server represents the HTTP server.
type Server struct {
	server *http.Server
}

// NewServer creates a new HTTP server.
func NewServer(ctx context.Context, cfg config.Config) *Server {
	tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))

	userHandler := handler.NewUserHandler(tmpl)

	// user routes
	router := handler.NewRouter()
	router.Mux.HandleFunc("GET /user/login", userHandler.Login)
	router.Mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// admin routes

	// middleware
	handlerMux := middleware.AccessLogWithCtx(ctx, router.Mux)

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:      handlerMux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		server: server,
	}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe(ctx context.Context) {
	const op = "server.ListenAndServe"
	err := s.server.ListenAndServe()
	if err != nil {
		logger.Z(ctx).Error(ctx, op, err.Error())
	}
}
