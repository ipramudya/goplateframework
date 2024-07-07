package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/worker/grpcserver"
	"github.com/goplateframework/pkg/db"
	"github.com/goplateframework/pkg/googlestorage"
	"github.com/goplateframework/pkg/logger"
	"google.golang.org/grpc"
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

	db, err := db.Init(conf)
	if err != nil {
		log.Fatalf("database connection error, %v", err)
		return err
	} else {
		log.Infof("database connected, status: %+v", db.Stats())
	}
	defer db.Close()

	// initialize google storage
	storage, err := googlestorage.Init(ctx, conf)

	if err != nil {
		log.Fatalf("google storage client error, %v", err)
		return err
	} else {
		log.Info("google storage client connected")
	}
	defer storage.Close()

	// set GRPC server up

	listener, err := net.Listen("tcp", conf.GRPCWorker.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	grpcServ := grpc.NewServer()
	grpcserver.Handle(grpcServ, &grpcserver.Options{
		Conf:    conf,
		DB:      db,
		Log:     log,
		Storage: storage,
	})

	// channel for storing grpc server errors which may occur during serving net listener
	serverErrCh := make(chan error, 1)

	// run grpc server in goroutine
	go func() {
		log.Infof("starting grpc server on port: %s", conf.GRPCWorker.Port)
		serverErrCh <- grpcServ.Serve(listener)
	}()

	select {
	case sig := <-shutdownCh:
		log.Infof("shutdown signal received: %v", sig)
		grpcServ.GracefulStop()
		defer func() {
			log.Infof("graceful stop completed: %s", sig)
			listener.Close()
		}()

	case err := <-serverErrCh:
		if err != nil {
			grpcServ.Stop()
			defer listener.Close()

			log.Fatalf("grpc server error, %v", err)
			return err
		}
	}

	return nil
}
