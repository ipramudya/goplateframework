package grpcserver

import (
	"cloud.google.com/go/storage"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/worker/pb"
	"github.com/goplateframework/pkg/logger"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type Options struct {
	Conf    *config.Config
	DB      *sqlx.DB
	Log     *logger.Log
	Storage *storage.Client
}

type server struct {
	pb.UnimplementedWorkerServer

	conf    *config.Config
	db      *sqlx.DB
	log     *logger.Log
	storage *storage.Client
}

func Handle(s *grpc.Server, opts *Options) {
	pb.RegisterWorkerServer(s, &server{
		conf:    opts.Conf,
		db:      opts.DB,
		log:     opts.Log,
		storage: opts.Storage,
	})
}
