package repository

import (
	"context"
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	"github.com/TestardR/geo-tracking/internal/domain/driver/model"
)

type CoordinatePersister interface {
	Persist(ctx context.Context, addCoordinateChange coordinateModel.AddCoordinateChange) error
}

type CoordinateFinder interface {
	Find(ctx context.Context, driverId model.DriverId) ([]coordinateModel.Coordinate, error)
}
