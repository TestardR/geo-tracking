package repository

import (
	"context"
	coordinateRepository "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	"github.com/TestardR/geo-tracking/internal/domain/driver/model"
)

type CoordinatePersister interface {
	Persist(ctx context.Context, addCoordinateChange coordinateRepository.AddCoordinateChange) error
}

type CoordinateFinder interface {
	Find(ctx context.Context, driverId model.DriverId) ([]coordinateRepository.Coordinate, error)
}
