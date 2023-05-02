package repository

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type StatusFinder interface {
	Find(ctx context.Context, driverId model.DriverId) (model.Status, error)
}

type StatusPersister interface {
	Persist(ctx context.Context, driverId model.DriverId, status model.Status) error
}
