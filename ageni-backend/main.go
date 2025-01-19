package main

import (
	"context"
	"flag"
	"github.com/readonme/open-studio/dal"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/readonme/open-studio/caller"
	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/conf"
	"github.com/readonme/open-studio/router"
)

func main() {
	flag.Parse()
	conf.Init()
	log.Init(conf.Conf.Log)
	dal.Init(conf.Conf.StudioDB)
	caller.InitCaller(conf.Conf)
	router.Init(conf.Conf)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	router.Stop(ctx)
	log.Info("Server Exited Properly")
}
