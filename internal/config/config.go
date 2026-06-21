package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
}

type AppConfig struct {
	Port int
	Env  string
}

type DBConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("app.port", 8080)
	v.SetDefault("app.env", "development")
	v.SetDefault("jwt.expires_in", "24h")

	if err := v.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
