package distance

import (
	"context"
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	distanceModel "github.com/TestardR/geo-tracking/internal/domain/distance"
)

const (
	HaversineFormula string = "haversine"
	VincentyFormula  string = "vincenty"
)

type StrategyExecutor interface {
	Distance(context.Context, []coordinateModel.Coordinate) (distanceModel.Distance, error)
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

func (d *distanceFinder) Distance(ctx context.Context, coordinates []coordinateModel.Coordinate) (distanceModel.Distance, error) {
	return d.strategies[d.strategy].Distance(ctx, coordinates)
}
