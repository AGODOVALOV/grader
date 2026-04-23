// Package client provides an HTTP server.
package client

import (
	"context"
	"embed"
	"html/template"
	"net"
	"net/http"
	"strconv"

	_ "github.com/AGODOVALOV/grader/docs/swagger"
	"github.com/AGODOVALOV/grader/pkg/client/middleware"
	"github.com/AGODOVALOV/grader/pkg/client/user"
	"github.com/AGODOVALOV/grader/pkg/client/user/repo"
	"github.com/AGODOVALOV/grader/pkg/client/user/usecase"
	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/AGODOVALOV/grader/pkg/token"
	"github.com/swaggo/http-swagger"
)

//go:embed html/templates/*.html
var templateFS embed.FS

// Server represents the HTTP server.
type Server struct {
	server *http.Server
	user   *user.User
	token  token.Maker
}

// NewClientServer creates a new HTTP server.
func NewClientServer(ctx context.Context, cfg *config.Config, repo *repo.Repo, fStorage *s3.FileStorage) (*Server, error) {
	tmpl := template.Must(template.ParseFS(templateFS, "html/templates/*.html"))

	tokenMaker, err := token.NewJWTMaker(&cfg.Token)
	if err != nil {
		return nil, err
	}

	usr := user.NewUser(tmpl, usecase.NewUserService(repo, fStorage, tokenMaker))

	// configure router
	r := configureRouter(usr)

	// add middleware
	handlerMux := middleware.AccessLogWithCtx(ctx, r)
	handlerMux = middleware.Auth(tokenMaker, handlerMux)

	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.WebServer.Host, strconv.Itoa(cfg.WebServer.Port)),
		Handler:      handlerMux,
		ReadTimeout:  cfg.WebServer.ReadTimeout,
		WriteTimeout: cfg.WebServer.WriteTimeout,
		IdleTimeout:  cfg.WebServer.IdleTimeout,
	}

	return &Server{
		server: srv,
		user:   usr,
		token:  tokenMaker,
	}, nil
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe(ctx context.Context) {
	const op = "server.ListenAndServe"
	err := s.server.ListenAndServe()
	if err != nil {
		logger.Z(ctx).Error(ctx, op, err.Error())
	}
}

func configureRouter(u *user.User) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /user/login", u.Handler.Login)
	router.HandleFunc("GET /admin", u.Handler.Admin)
	router.HandleFunc("GET /user/register", u.Handler.Register)
	router.HandleFunc("GET /user/account/{userID}", u.Handler.Account)

	router.HandleFunc("POST /user/create", u.Handler.CreateUser)
	router.HandleFunc("POST /user/login", u.Handler.LoginUser)

	router.HandleFunc("POST /task/review", u.Handler.UploadTask)

	router.Handle("/swagger/", httpSwagger.WrapHandler)

	return router
}
