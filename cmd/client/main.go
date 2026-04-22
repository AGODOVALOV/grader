// Package main contains the main function for the client server.
package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	server "github.com/AGODOVALOV/grader/pkg/client"
	"github.com/AGODOVALOV/grader/pkg/client/user/repo"
	"github.com/AGODOVALOV/grader/pkg/config"
	"github.com/AGODOVALOV/grader/pkg/logger"
)

// main initializes the application by loading configuration, setting up the logger, and starting the client server.
func main() {
	// create and read config
	appCfg, err := config.GetApplicationConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// create logger
	z, err := logger.NewAppLogger(appCfg.GetConfig().Log)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	// ctx add logger
	ctx = logger.CtxWithLogger(ctx, z)
	defer cancel()

	// start client server
	z.Info(
		ctx,
		"client server",
		"starting",
		map[string]string{
			"url": net.JoinHostPort(appCfg.GetConfig().WebServer.Host, strconv.Itoa(appCfg.GetConfig().WebServer.Port)),
		})

	repoDB, err := repo.NewRepo(ctx, appCfg.GetConfig())
	if err != nil {
		z.Error(ctx, "DB connection", err.Error())
		return
	}

	srv, err := server.NewClientServer(ctx, appCfg.GetConfig(), repoDB)
	if err != nil {
		z.Error(ctx, "create server", err.Error())
		return
	}

	srv.ListenAndServe(ctx)
}
