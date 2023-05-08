package distance

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/domain/model"
)

const (
	HaversineFormula string = "haversine"
	VincentyFormula  string = "vincenty"
)

type StrategyExecutor interface {
	Distance(context.Context, []model.Coordinate) (float64, error)
}

type Strategy string

type distanceFinder struct {
	strategy   Strategy
	strategies map[Strategy]StrategyExecutor
}

func NewDistanceFinder(
	strategy Strategy,
	strategies map[Strategy]StrategyExecutor,
) *distanceFinder {
	return &distanceFinder{strategy: strategy, strategies: strategies}
}

func (d *distanceFinder) Distance(ctx context.Context, coordinates []model.Coordinate) (float64, error) {
	return d.strategies[d.strategy].Distance(ctx, coordinates)
}
