package config

type RpcClientCfg struct {
	Servers        []string `hcl:"servers" json:"servers" yaml:"servers"`
	ConnectTimeout int      `hcl:"connect-timeout" json:"connect-timeout" yaml:"connect-timeout"`
	Retry          int      `hcl:"retry" json:"retry" yaml:"retry"`
}
