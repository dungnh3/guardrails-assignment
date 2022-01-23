package config

import (
	"bytes"
	"strings"

	"github.com/dungnh3/guardrails-assignment/pkg/database"
	"github.com/dungnh3/guardrails-assignment/pkg/log"
	"github.com/dungnh3/guardrails-assignment/pkg/server"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type DeployEnv string

const (
	DeployDevelopmentEnv DeployEnv = "development"
	DeployProductionEnv            = "production"
)

type Config struct {
	Env        DeployEnv                 `json:"env,omitempty"`
	Logger     log.Config                `json:"logger"`
	PostgreSQL database.PostgreSQLConfig `json:"postgresql"`
	Server     server.Config             `json:"server"`
}

// loadDefaultConfig return a default object configuration
func loadDefaultConfig() *Config {
	return &Config{
		Env:        DeployProductionEnv,
		Logger:     log.DefaultConfig(),
		Server:     server.DefaultConfig(),
		PostgreSQL: database.PostgresSQLDefaultConfig(),
	}
}

func Load() (*Config, error) {
	// You should set default config value here
	c := loadDefaultConfig()
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	var err error
	if err = viper.ReadConfig(bytes.NewBuffer([]byte(defaultValue))); err != nil {
		return nil, err
	}

	err = viper.Unmarshal(c, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	})
	return c, err
}
