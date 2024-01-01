package config

type Log struct {
	// Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	Level            int    `default:"1" yaml:"level" help:"the application log level"`
	TimeFormat       string `default:"2006-01-02T15:04:05Z07:00" yaml:"time-format" help:"the application log time format"`
	TimestampEnabled bool   `yaml:"timestamp-enabled" default:"false"`
	ServiceName      string `yaml:"service-name" help:"the application service name"`
}
