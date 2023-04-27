package config

import (
	"github.com/caarlos0/env/v6"
)

type Application struct {
	DistanceAlgorithm string `env:"DISTANCE_ALGORITHM" example:"haversine|vincenty" envDefault:"haversine"`
}

func NewApplication() (*Application, error) {
	appConfig := &Application{}

	if err := env.Parse(appConfig); err != nil {
		return nil, err
	}

	return appConfig, nil
}
