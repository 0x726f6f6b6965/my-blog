package config

type Grpc struct {
	Host       string `yaml:"host" help:"the host to bind for GRPC server"`
	Port       int    `yaml:"port" help:"the port to bind for GRPC server"`
	Network    string `yaml:"network" default:"tcp" help:"the network type for GRPC server"`
	Socketpath string `yaml:"socketpath" help:"the socket path when network is unix"`
}

type Rest struct {
	Host string `yaml:"host" help:"the host to bind for REST server"`
	Port int    `yaml:"port" help:"the port to bind for REST server"`
}
