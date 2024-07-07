package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/httpserver"
	"github.com/goplateframework/internal/worker/pb"
	"github.com/goplateframework/pkg/db"
	"github.com/goplateframework/pkg/grpcclient"
	"github.com/goplateframework/pkg/logger"
	"github.com/goplateframework/pkg/redisdb"
)

func main() {
	// load and read config file
	rawConfig, err := config.Load()
	if err != nil {
		panic("load config error, " + err.Error())
	}

	// parse raw config into config struct
	conf, err := config.Parse(rawConfig)
	if err != nil {
		panic("parse config error, " + err.Error())
	}

	log := logger.Init(conf)

	ctx := context.Background()

	// start server in background
	if err := run(ctx, conf, log); err != nil {
		log.Fatalf("run error, %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, conf *config.Config, log *logger.Log) error {
	log.Infof("starting server...")

	// retrieve database connection
	log.Infof("initializing database connection on host: %s", conf.DB.Host)

	db, err := db.Init(conf)
	if err != nil {
		log.Fatalf("database connection error, %v", err)
		return err
	} else {
		log.Infof("database connected, status: %+v", db.Stats())
	}
	defer db.Close()

	// retrieve redis connection

	rdb, err := redisdb.Init(conf)
	if err != nil {
		log.Fatalf("redis connection error, %v", err)
		return err
	} else {
		log.Infof("redis connected, status: %+v", rdb.PoolStats())
	}
	defer rdb.Close()

	// initialize grpc client caller

	grpcconn, err := grpcclient.Init(conf)
	if err != nil {
		log.Fatalf("grpc client error, %v", err)
		return err
	} else {
		log.Infof("grpc client connected on %s%s", conf.Server.Host, conf.GRPCWorker.Port)
	}
	defer grpcconn.Close()

	// channel to receive shutdownCh signal, for graceful shutdownCh
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// initialize http server by passing necessary dependencies
	server := httpserver.Init(&httpserver.Options{
		DB:       db,
		Cache:    rdb,
		Log:      log,
		ServConf: conf,
		Worker:   pb.NewWorkerClient(grpcconn),
	})

	// channel for handling server errors which may occur during listening and serving
	serverErrCh := make(chan error, 1)

	// run server in goroutine
	go func() {
		// create http server, pass server configuration to echo instance
		s := &http.Server{
			Addr:         conf.Server.Host + ":" + conf.Server.Port,
			ReadTimeout:  time.Second * conf.Server.ReadTimeout,
			WriteTimeout: time.Second * conf.Server.WriteTimeout,
		}

		log.Infof("server started on %s:%s", conf.Server.Host, conf.Server.Port)
		serverErrCh <- server.StartServer(s)
	}()

	// main thread waits for shutdown signal within graceful shutdown
	// or server error channel, and then handle it accordingly
	select {
	case sig := <-shutdownCh:
		log.Infof("shutdown started: %s", sig)
		defer log.Infof("shutdown completed: %s", sig)

		ctx, cancel := context.WithTimeout(ctx, conf.Server.CtxDefaultTimeout*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("gracefull shutdown failed, server forced to shutdown: %v", err)
		}

	case err := <-serverErrCh:
		log.Errorf("server error, %v", err)
		return err
	}

	return nil
}
