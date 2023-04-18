package config

import (
	"github.com/jessevdk/go-flags"
)

type Config struct {
	Env      string `long:"env" env:"PIM_ENV" example:"live|dev|staging" default:"dev"`
	LogLevel string `long:"log-level" env:"GT_LOG_LEVEL" example:"debug,warn" default:"error"`
	LogPath  string `long:"log-path" env:"GT_LOG_PATH" example:"log/file.log" default:""`

	HttpPort string `long:"http-port" env:"PIM_HTTP_PORT" example:":8091" default:":80"`

	RedisMasterAddr string `long:"redis-addr" env:"GT_REDIS_ADDR" example:"localhost:6379" default:"localhost:6379"`
	RedisDb         int    `long:"redis-db" env:"GT_REDIS_DB" example:"0" default:"0"`

	NatsBrokerList string `long:"nats-broker-list" env:"NATS_BROKER_LIST" example:"localhost:4222" default:"localhost:4222"`
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
