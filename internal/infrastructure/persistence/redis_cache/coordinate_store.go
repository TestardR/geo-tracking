package redis_cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/infrastructure/persistence/redis_cache/entity"
)

type coordinateStore struct {
	redis *client
}

func NewCoordinateStore(redis *client) *coordinateStore {
	return &coordinateStore{redis: redis}
}

func (s *coordinateStore) Persist(
	ctx context.Context,
	driverId model.DriverId,
	coordinate model.Coordinate,
) error {
	coordinatesAsBytes, err := s.redis.Get(ctx, driverId.Id())
	if err != nil {
		return err
	}

	var coordinates []entity.Coordinate
	err = json.Unmarshal(coordinatesAsBytes, &coordinates)
	if err != nil {
		return err
	}

	coordinates = coordinates[:len(coordinates)-1]

	c := entity.Coordinate{
		Longitude: coordinate.Longitude(),
		Latitude:  coordinate.Latitude(),
	}
	coordinates = append(coordinates, c)

	coordinatesAsBytes, err = json.Marshal(coordinates)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, driverId.Id(), coordinatesAsBytes, time.Duration(0))
}
