package config

import "time"

type Client struct {
	Host              string        `yaml:"host" help:"the host to connect server"`
	GrpcPort          int           `yaml:"grpc-port" help:"the port to bind for GRPC server"`
	GrpcNetwork       string        `yaml:"grpc-network" default:"tcp" help:"the network type for GRPC server"`
	GrpcSocketpath    string        `yaml:"grpc-socketpath" help:"the socket path when network is unix"`
	HTTPPort          int           `yaml:"port" help:"the port to bind for REST server"`
	ConnectionTimeout time.Duration `yaml:"connection-timeout" help:"the service connection timeout"`
}
