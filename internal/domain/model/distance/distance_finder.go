package distance

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type StrategyExecutor interface {
	Distance(context.Context, []model.Coordinate) (float64, error)
}

type Strategy string

type distanceFinder struct {
	strategy   Strategy
	strategies map[Strategy]StrategyExecutor
}

func (d *distanceFinder) Distance(ctx context.Context, coordinates []model.Coordinate) (float64, error) {
	return d.strategies[d.strategy].Distance(ctx, coordinates)
}
