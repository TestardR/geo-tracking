package redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/persistence/redis_cache/entity"
)

const maxNumberOfCoordinates = 30

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
	coordinatesAsBytes, err := s.redis.Get(ctx, fmt.Sprintf("c:%s", driverId.Id()))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			var coordinates []entity.Coordinate
			coordinates = append(coordinates, coordinateModelToEntity(coordinate))

			coordinatesAsBytes, err = json.Marshal(coordinates)
			if err != nil {
				return err
			}

			return s.redis.Set(ctx, driverId.Id(), coordinatesAsBytes, time.Duration(0))
		}
		return err
	}

	var coordinates []entity.Coordinate
	err = json.Unmarshal(coordinatesAsBytes, &coordinates)
	if err != nil {
		return err
	}

	if len(coordinates) >= maxNumberOfCoordinates {
		coordinates = coordinates[:len(coordinates)-1]
	}

	coordinates = coordinates[:len(coordinates)-1]
	coordinates = append(coordinates, coordinateModelToEntity(coordinate))

	coordinatesAsBytes, err = json.Marshal(coordinates)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, driverId.Id(), coordinatesAsBytes, time.Duration(0))
}

func (s *coordinateStore) Find(
	ctx context.Context,
	driverId model.DriverId,
) ([]model.Coordinate, error) {
	coordinatesAsBytes, err := s.redis.Get(ctx, driverId.Id())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, shared.NewDomainError(repository.ErrDriverIdNotFoundMessage)
		}
		return nil, err
	}

	var coordinates []entity.Coordinate
	err = json.Unmarshal(coordinatesAsBytes, &coordinates)
	if err != nil {
		return nil, err
	}

	var modelCoordinates []model.Coordinate
	for _, coordinate := range coordinates {
		modelCoordinates = append(modelCoordinates,
			model.RecreateCoordinate(coordinate.Longitude, coordinate.Latitude),
		)
	}

	return modelCoordinates, nil
}

func coordinateModelToEntity(coordinate model.Coordinate) entity.Coordinate {
	return entity.Coordinate{
		Longitude: coordinate.Longitude(),
		Latitude:  coordinate.Latitude(),
	}
}
