package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/storage"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/worker/grpcserver"
	"github.com/goplateframework/pkg/db"
	"github.com/goplateframework/pkg/logger"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func main() {
	// load and read config file
	rawConfig, err := config.LoadConfig()
	if err != nil {
		panic("load config error, " + err.Error())
	}

	// parse raw config into config struct
	conf, err := config.ParseConfig(rawConfig)
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

	log.Infof("initializing database connection on host: %s", conf.DB.Host)

	// retrieve database connection
	db, err := db.Init(conf)
	if err != nil {
		log.Fatalf("database connection error, %v", err)
		return err
	} else {
		log.Infof("database connected, status: %+v", db.Stats())
	}
	defer db.Close()

	log.Info("initializing firebase storage client...")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current working directory error, %v", err)
	}

	// initialize firebase storage
	storagePath := cwd + conf.Firebase.Path
	storage, err := storage.NewClient(ctx, option.WithCredentialsFile(storagePath))
	if err != nil {
		log.Fatalf("firebase storage client error, %v", err)
		return err
	} else {
		log.Info("firebase storage client connected")
	}
	defer storage.Close()

	// set up GRPC server
	lis, err := net.Listen("tcp", conf.RPC.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	grpcSer := grpc.NewServer()
	grpcserver.Handle(grpcSer, &grpcserver.Options{
		Conf:    conf,
		DB:      db,
		Log:     log,
		Storage: storage,
	})

	serverErrCh := make(chan error, 1)
	go func() {
		log.Infof("starting grpc server on port: %s", conf.RPC.Port)
		serverErrCh <- grpcSer.Serve(lis)
	}()

	select {
	case sig := <-shutdownCh:
		log.Infof("shutdown signal received: %v", sig)
		grpcSer.GracefulStop()
		defer func() {
			log.Infof("graceful stop completed: %s", sig)
			lis.Close()
		}()

	case err := <-serverErrCh:
		if err != nil {
			grpcSer.Stop()
			defer lis.Close()

			log.Fatalf("grpc server error, %v", err)
			return err
		}
	}

	return nil
}
