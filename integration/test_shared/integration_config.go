package test_shared

import (
	"sync"
	"testing"

	"github.com/caarlos0/env/v6"
)

var configLoader sync.Once
var integrationConfig *IntegrationConfig

type IntegrationConfig struct {
	Env      string `long:"env" env:"GT_ENV" example:"live|dev|staging" default:"dev"`
	LogLevel string `long:"log-level" env:"GT_LOG_LEVEL" example:"debug,warn" default:"error"`
	LogPath  string `long:"log-path" env:"GT_LOG_PATH" example:"log/file.log" default:""`

	HttpPort string `long:"http-port" env:"GT_HTTP_PORT" example:":8090" default:":8090"`

	RedisMasterAddr string `long:"redis-addr" env:"GT_REDIS_ADDR" example:"localhost:6379" default:"localhost:6379"`
	RedisDb         int    `long:"redis-db" env:"GT_REDIS_DB" example:"0" default:"0"`

	NatsBrokerList string `long:"nats-broker-list" env:"NATS_BROKER_LIST" example:"localhost:4222" default:"localhost:4222"`
}

func GetIntegrationConfig(t *testing.T) *IntegrationConfig {
	if integrationConfig == nil {
		integrationConfig = &IntegrationConfig{}
		configLoader.Do(func() {
			if err := env.Parse(integrationConfig); err != nil {
				t.Fatal(err)
			}
		})
	}

	return integrationConfig
}
