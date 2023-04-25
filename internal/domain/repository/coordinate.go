package repository

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type CoordinatePersister interface {
	Persist(ctx context.Context, driverId model.DriverId, coordinate model.Coordinate) error
}

type CoordinateFinder interface {
	Find(ctx context.Context, driverId model.DriverId) ([]model.Coordinate, error)
}
