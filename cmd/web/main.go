// Package main contains the main function for the web server.
package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/AGODOVALOV/grader/pkg/config"
	"github.com/AGODOVALOV/grader/pkg/logger"
	server "github.com/AGODOVALOV/grader/pkg/webserver"
)

// main initializes the application by loading configuration, setting up the logger, and starting the web server.
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

	// start web server
	z.Info(
		ctx,
		"web server",
		"starting",
		map[string]string{
			"url": net.JoinHostPort(appCfg.GetConfig().WebServer.Host, strconv.Itoa(appCfg.GetConfig().WebServer.Port)),
		})
	server.NewServer(ctx, appCfg.GetConfig().WebServer).ListenAndServe(ctx)
}
