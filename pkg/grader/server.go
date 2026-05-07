package grader

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/grader/grader"
	"github.com/AGODOVALOV/grader/pkg/grader/middleware"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/token"
)

// Server represents a server.go server instance with an HTTP server and token management capabilities.
type Server struct {
	server *http.Server
	token  token.Maker
	Grader *grader.Grader
}

// NewGraderServer creates a new HTTP server.
func NewGraderServer(ctx context.Context,
	cfg *config.Config,
	graderProc *grader.Grader,
	tokenMaker token.Maker) (*Server, error) {
	// configure router
	router := configureRouter(graderProc)

	// add middleware
	handlerMux := middleware.AccessLogWithCtx(ctx, router)
	handlerMux = middleware.Auth(tokenMaker, handlerMux)

	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Grader.Server.Host, strconv.Itoa(cfg.Grader.Server.Port)),
		Handler:      handlerMux,
		ReadTimeout:  cfg.Grader.Server.ReadTimeout,
		WriteTimeout: cfg.Grader.Server.WriteTimeout,
		IdleTimeout:  cfg.Grader.Server.IdleTimeout,
	}

	return &Server{
		server: srv,
		token:  tokenMaker,
		Grader: graderProc,
	}, nil
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe(ctx context.Context) {
	const op = "server.ListenAndServe"

	defer func() {
		close(s.Grader.Handler.GraderService.WP.Tasks)
	}()

	err := s.server.ListenAndServe()
	if err != nil {
		logger.Z(ctx).Error(ctx, op, err.Error())
	}
}

func configureRouter(g *grader.Grader) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/grader", g.Handler.Grade)
	return router
}
