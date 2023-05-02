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
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/redis_cache/entity"
	redisCache "github.com/TestardR/geo-tracking/internal/infrastructure/shared/redis_cache"
)

type coordinateStore struct {
	redis *redisCache.Client
}

func NewCoordinateStore(redis *redisCache.Client) *coordinateStore {
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
			coordinates = append(coordinates, CoordinateModelToEntity(coordinate))

			coordinatesAsBytes, err = json.Marshal(coordinates)
			if err != nil {
				return err
			}

			return s.redis.Set(ctx, fmt.Sprintf("c:%s", driverId.Id()), coordinatesAsBytes, time.Duration(0))
		}

		return err
	}

	var coordinates []entity.Coordinate
	err = json.Unmarshal(coordinatesAsBytes, &coordinates)
	if err != nil {
		return err
	}

	inTimeWindowCoordinates := make([]entity.Coordinate, 0, len(coordinates))
	t := time.Now()
	for _, c := range coordinates {
		if t.Sub(c.CreatedAt).Minutes() < 5 {
			inTimeWindowCoordinates = append(inTimeWindowCoordinates, c)
		}

	}

	inTimeWindowCoordinates = append(inTimeWindowCoordinates, CoordinateModelToEntity(coordinate))
	coordinatesAsBytes, err = json.Marshal(inTimeWindowCoordinates)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, fmt.Sprintf("c:%s", driverId.Id()), coordinatesAsBytes, time.Duration(0))
}

func (s *coordinateStore) Find(
	ctx context.Context,
	driverId model.DriverId,
) ([]model.Coordinate, error) {
	coordinatesAsBytes, err := s.redis.Get(ctx, fmt.Sprintf("c:%s", driverId.Id()))
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
			model.RecreateCoordinate(
				coordinate.Longitude,
				coordinate.Latitude,
				coordinate.CreatedAt,
			),
		)
	}

	return modelCoordinates, nil
}

func CoordinateModelToEntity(coordinate model.Coordinate) entity.Coordinate {
	return entity.Coordinate{
		Longitude: coordinate.Longitude(),
		Latitude:  coordinate.Latitude(),
		CreatedAt: coordinate.CreatedAt(),
	}
}
