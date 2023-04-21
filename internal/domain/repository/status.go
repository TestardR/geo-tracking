package repository

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/infrastructure/persistence/redis_cache/entity"
)

type FindStatus interface {
	// TODO: should be model.DriverId
	Find(ctx context.Context, driver entity.Driver) (model.Status, error)
}
