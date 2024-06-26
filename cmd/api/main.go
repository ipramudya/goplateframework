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
	"github.com/goplateframework/internal/server"
	"github.com/goplateframework/pkg/db"
	"github.com/goplateframework/pkg/logger"
	"github.com/goplateframework/pkg/redisdb"
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

	log.Info("initializing redis connection...")

	// retrieve redis connection
	rdb, err := redisdb.Init(conf)
	if err != nil {
		log.Fatalf("redis connection error, %v", err)
		return err
	} else {
		log.Info("redis connected")
	}
	defer rdb.Close()

	log.Infof("starting server on port: %s", conf.Server.Port)

	// channel to receive shutdown signal, for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// server configuration
	serverConf := server.Config{
		DB:       db,
		RDB:      rdb,
		Log:      log,
		ServConf: conf,
	}

	// create http server, pass server configuration to server handler
	serv := &http.Server{
		Addr:         conf.Server.Host + ":" + conf.Server.Port,
		Handler:      server.Handler(&serverConf),
		ReadTimeout:  time.Second * conf.Server.ReadTimeout,
		WriteTimeout: time.Second * conf.Server.WriteTimeout,
	}

	// channel for handling server errors which may occur during listening and serving
	serverErrs := make(chan error, 1)

	// run server in goroutine
	go func() {
		log.Infof("server successfully running on %s", serv.Addr)
		serverErrs <- serv.ListenAndServe()
	}()

	// main thread waits for shutdown signal within graceful shutdown
	// or server error channel, and then handle it accordingly
	select {
	case sig := <-shutdown:
		log.Infof("shutdown started: %s", sig)
		defer log.Infof("shutdown completed: %s", sig)

		ctx, cancel := context.WithTimeout(ctx, conf.Server.CtxDefaultTimeout*time.Second)
		defer cancel()

		if err := serv.Shutdown(ctx); err != nil {
			serv.Close()
			return fmt.Errorf("gracefull shutdown failed, server forced to shutdown: %v", err)
		}

	case err := <-serverErrs:
		log.Errorf("server error, %v", err)
		return err
	}

	return nil
}
