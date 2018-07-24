package server

import (
	"net"

	"../api"
	"google.golang.org/grpc"
)

func Serve(addr string, root string, sizeLimit int, logPath string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s, err := api.NewServer(root, sizeLimit, logPath)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	api.RegisterImageServiceServer(grpcServer, s)
	return grpcServer.Serve(l)
}
