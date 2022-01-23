package server

import (
	"fmt"
)

// Config hold http/grpc server config
type Config struct {
	GRPC Listen `json:"grpc" mapstructure:"grpc" yaml:"grpc"`
	HTTP Listen `json:"http" mapstructure:"http" yaml:"http"`
}

// Listen config for host/port socket listener
type Listen struct {
	Host string `json:"host" mapstructure:"host" yaml:"host"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}

// DefaultConfig return default server config
func DefaultConfig() Config {
	return Config{
		GRPC: Listen{
			Host: "0.0.0.0",
			Port: 10443,
		},
		HTTP: Listen{
			Host: "0.0.0.0",
			Port: 10080,
		},
	}
}

// String return socket listen data source name
func (l Listen) String() string {
	return fmt.Sprintf("%v:%v", l.Host, l.Port)
}
