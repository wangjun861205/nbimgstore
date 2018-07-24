package client

import (
	"../api"
	"google.golang.org/grpc"
)

type Client struct {
	api.ImageServiceClient
	conn *grpc.ClientConn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{api.NewImageServiceClient(conn), conn}, nil
}
