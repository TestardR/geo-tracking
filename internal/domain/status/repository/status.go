package repository

import (
	"context"
	driverModel "github.com/TestardR/geo-tracking/internal/domain/driver/model"
	statusModel "github.com/TestardR/geo-tracking/internal/domain/status/model"
)

type StatusFinder interface {
	Find(ctx context.Context, driverId driverModel.DriverId) (statusModel.Status, error)
}

type StatusPersister interface {
	Persist(ctx context.Context, driverId driverModel.DriverId, status statusModel.Status) error
}
