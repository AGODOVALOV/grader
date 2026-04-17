// Package client provides an HTTP server.
package client

import (
	"context"
	"embed"
	"html/template"
	"net"
	"net/http"
	"strconv"

	"github.com/AGODOVALOV/grader/pkg/client/config"
	handler2 "github.com/AGODOVALOV/grader/pkg/client/user/handler"
	"github.com/AGODOVALOV/grader/pkg/client/user/middleware"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/swaggo/http-swagger"
)

//go:embed html/templates/*.html
var templateFS embed.FS

// Server represents the HTTP server.
type Server struct {
	server *http.Server
}

// NewClientServer creates a new HTTP server.
func NewClientServer(ctx context.Context, cfg config.Config) *Server {
	tmpl := template.Must(template.ParseFS(templateFS, "html/templates/*.html"))

	userHandler := handler2.NewUserHandler(tmpl)

	// user routes
	router := http.NewServeMux()
	router.HandleFunc("GET /user/login", userHandler.Login)
	router.HandleFunc("GET /user/register", userHandler.Register)
	router.HandleFunc("GET /user/account/{userID}", userHandler.Account)
	router.Handle("/swagger/", httpSwagger.WrapHandler)

	// admin routes

	// middleware
	handlerMux := middleware.AccessLogWithCtx(ctx, router)

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
