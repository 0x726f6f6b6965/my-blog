package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	conn *grpc.ClientConn
}

func NewGrpcConn(ctx context.Context, uri string, timeout time.Duration) (*GrpcClient, error) {
	conn, err := grpc.DialContext(ctx, uri,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithIdleTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return &GrpcClient{conn: conn}, nil
}

func (c *GrpcClient) Reconnect() {
	if c.conn.GetState() != connectivity.Ready {
		c.conn.Connect()
	}
}

func (c *GrpcClient) GetConn() *grpc.ClientConn {
	return c.conn
}
