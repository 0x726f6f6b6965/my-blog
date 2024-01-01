package config

type RedisCfg struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	MaxRetries int    `yaml:"max-retries" default:"3"`
	DB         int    `yaml:"db"`
}
