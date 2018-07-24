package server

import (
	"log"
	"net"

	"../api"
	"google.golang.org/grpc"
)

func NewServer(addr string, root string, sizeLimit int, logPath string) (*grpc.Server, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s, err := api.NewServer(root, sizeLimit, logPath)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer()
	api.RegisterImageServiceServer(grpcServer, s)
	go func() {
		log.Fatal(grpcServer.Serve(l))
	}()
	return grpcServer, nil
}
