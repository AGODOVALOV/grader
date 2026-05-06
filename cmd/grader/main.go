// Package main contains the main function for the server.go server.
package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/AGODOVALOV/grader/pkg/config"
	graderserver "github.com/AGODOVALOV/grader/pkg/grader"
	"github.com/AGODOVALOV/grader/pkg/grader/grader"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/AGODOVALOV/grader/pkg/token"
)

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
			"url": net.JoinHostPort(appCfg.GetConfig().Grader.Server.Host, strconv.Itoa(appCfg.GetConfig().Grader.Server.Port)),
		})

	// init file storage
	fStorage, err := s3.NewFileStorage(ctx, &appCfg.GetConfig().FileStorage)
	if err != nil {
		z.Error(ctx, "init file storage", err.Error())
		return
	}

	tokenMaker, err := token.NewJWTMaker(&appCfg.GetConfig().Token)
	if err != nil {
		z.Error(ctx, "init token maker", err.Error())
		return
	}

	// init grader processor
	graderProc := grader.NewGrader(ctx, fStorage, tokenMaker, appCfg.GetConfig())

	// init web server
	srv, err := graderserver.NewGraderServer(ctx,
		appCfg.GetConfig(),
		graderProc,
		tokenMaker,
		fStorage)
	if err != nil {
		z.Error(ctx, "create server", err.Error())
		return
	}

	//start worker pool
	graderProc.Handler.GraderService.WP.StartProcessingGradeTasks(ctx)

	// start server
	srv.ListenAndServe(ctx)

	// close channel
	close(graderProc.Handler.GraderService.WP.Tasks)

	//wait for all tasks to finish
	graderProc.Handler.GraderService.WP.Wait()

}
