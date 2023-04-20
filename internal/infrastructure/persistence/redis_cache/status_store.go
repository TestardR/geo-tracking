package redis_cache

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type statusStore struct {
	redis *client
	// lib to compute distance between two coordinates
}

func NewStatusStore(redis *client) *statusStore {
	return &statusStore{redis: redis}
}

func (s *statusStore) Find(
	ctx context.Context,
	driverId model.DriverId,
) (model.Status, error) {
	return model.Status{}, nil
}
