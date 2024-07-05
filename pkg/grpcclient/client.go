package grpcclient

import (
	"github.com/goplateframework/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init(conf *config.Config) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		conf.RPC.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
