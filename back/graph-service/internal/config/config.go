package config

import (
	"time"

	libCfg "github.com/0x726f6f6b6965/my-blog/lib/config"
)

type Config struct {
	Name    string          `yaml:"name" help:"the application name"`
	Rest    libCfg.Rest     `yaml:"rest"`
	Redis   libCfg.RedisCfg `yaml:"redis" help:"the application redis option"`
	Clients struct {
		Blog   Client `yaml:"blog-svc"`
		User   Client `yaml:"user-svc"`
		Search Client `yaml:"search-svc"`
	} `yaml:"clients"`
	Log libCfg.Log `yaml:"log" help:"the application log"`
}

type Client struct {
	Grpc              libCfg.Grpc   `yaml:"grpc" help:"the application grpc option"`
	ConnectionTimeout time.Duration `yaml:"connection-timeout"`
}
