// Package main contains the main function for the client server.
package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	server "github.com/AGODOVALOV/grader/pkg/client"
	"github.com/AGODOVALOV/grader/pkg/client/metrics"
	"github.com/AGODOVALOV/grader/pkg/client/outbox"
	"github.com/AGODOVALOV/grader/pkg/client/user/repo"
	"github.com/AGODOVALOV/grader/pkg/config"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
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

	// init db
	repoDB, err := repo.NewRepo(ctx, appCfg.GetConfig())
	if err != nil {
		z.Error(ctx, "DB connection", err.Error())
		return
	}

	// init file storage
	fStorage, err := s3.NewFileStorage(ctx, &appCfg.GetConfig().FileStorage)
	if err != nil {
		z.Error(ctx, "init file storage", err.Error())
		return
	}

	//init metrics collector
	metricsCollector := metrics.NewCollector(ctx)

	// init web server
	srv, err := server.NewClientServer(ctx, appCfg.GetConfig(), repoDB, fStorage, metricsCollector)
	if err != nil {
		z.Error(ctx, "create server", err.Error())
		return
	}

	// start outbox transfer
	outboxTransfer, err := outbox.NewOutbox(ctx, srv.User, &appCfg.GetConfig().MsgQueue)
	if err != nil {
		z.Error(ctx, "starting outbox process", err.Error())
		return
	}

	go func() {
		err = outboxTransfer.StartSending(ctx)
		if err != nil {
			z.Error(ctx, "start outbox transfer", err.Error())
			return
		}
	}()

	// start server
	srv.ListenAndServe(ctx)
}
