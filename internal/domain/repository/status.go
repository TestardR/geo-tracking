package repository

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type StatusFinder interface {
	Find(ctx context.Context, driverId model.DriverId) (model.Status, error)
}
