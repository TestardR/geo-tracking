package config

import (
	"github.com/jessevdk/go-flags"
)

type Config struct {
	LogLevel string `long:"log-level" env:"PIM_LOG_LEVEL" example:"debug,warn" default:"error"`

	HttpPort string `long:"http-port" env:"PIM_HTTP_PORT" example:":8091" default:":80"`

	RedisMasterAddr string `long:"redis-addr" env:"GT_REDIS_ADDR" example:"localhost:6379" default:"localhost:6379"`
	RedisDb         int    `long:"redis-db" env:"GT_REDIS_DB" example:"0" default:"0"`
}

func NewConfig() (*Config, error) {
	var config = &Config{}

	parser := flags.NewParser(config, flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		return config, err
	}

	return config, nil
}
