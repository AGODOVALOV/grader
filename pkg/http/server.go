package http

import (
	"context"
	"embed"
	"html/template"
	"net"
	"net/http"
	"strconv"

	_ "github.com/AGODOVALOV/grader/docs"
	"github.com/AGODOVALOV/grader/pkg/http/config"
	"github.com/AGODOVALOV/grader/pkg/http/handler"
	"github.com/AGODOVALOV/grader/pkg/http/handler/user"
	"github.com/AGODOVALOV/grader/pkg/http/middleware"
	"github.com/AGODOVALOV/grader/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed templates/*.html
var templateFS embed.FS

type Server struct {
	server *http.Server
}

func NewServer(ctx context.Context, cfg config.Config) *Server {

	tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))

	userHandler := user.NewUserHandler(tmpl)

	// user routes
	router := handler.NewRouter()
	router.Mux.HandleFunc("GET /user/login", userHandler.Login)
	router.Mux.Handle("/swagger/", httpSwagger.WrapHandler)

	//admin routes

	//middleware
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

func (s *Server) ListenAndServe(ctx context.Context) {
	const op = "server.ListenAndServe"
	err := s.server.ListenAndServe()
	if err != nil {
		logger.Z(ctx).Error(ctx, op, err.Error())
	}
}
