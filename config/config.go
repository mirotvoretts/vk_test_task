package config

import "time"

type Config struct {
	GRPCPort        int           `mapstructure:"GRPC_PORT"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

func New() *Config {
	return &Config{
		GRPCPort:        50051,
		ShutdownTimeout: 10 * time.Second,
	}
}
