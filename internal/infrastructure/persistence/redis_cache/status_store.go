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

type coordinateFinder interface {
	Find(context.Context, entity.Driver) ([]model.Coordinate, error)
}

type distanceFinder interface {
	Find(context.Context, []model.Coordinate) float64
}

type statusStore struct {
	redis           *client
	coordinateStore coordinateFinder
	distanceFinder  distanceFinder
}

func NewStatusStore(
	redis *client,
	coordinateStore coordinateFinder,
	distanceFinder distanceFinder,
) *statusStore {
	return &statusStore{
		redis:           redis,
		coordinateStore: coordinateStore,
		distanceFinder:  distanceFinder,
	}
}

func (s *statusStore) Find(ctx context.Context, driver entity.Driver) (model.Status, error) {
	driverKey := fmt.Sprintf("s:%s", driver.Id)
	statusAsBytes, err := s.redis.Get(ctx, driverKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Status{}, shared.NewDomainError(repository.ErrDriverIdNotFoundMessage)
		}

		return model.Status{}, err
	}

	if len(statusAsBytes) > 0 {
		var status entity.Status
		err = json.Unmarshal(statusAsBytes, &status)
		if err != nil {
			return model.Status{}, err
		}

		return s.toDomain(status), nil
	}

	coordinates, err := s.coordinateStore.Find(ctx, driver)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// return 404
			return model.Status{}, err
		}

		// 500
		return model.Status{}, err
	}

	status := model.NewStatus(false)
	status.ComputeZombieStatus(s.distanceFinder.Find(ctx, coordinates))

	entityStatus := s.toEntity(status)
	entityStatusAsBytes, err := json.Marshal(entityStatus)
	if err != nil {
		// 500
		return model.Status{}, err
	}
	err = s.redis.Set(ctx, driverKey, entityStatusAsBytes, time.Duration(0))
	if err != nil {
		// 500
		return model.Status{}, err
	}

	return status, nil
}

func (s *statusStore) toDomain(status entity.Status) model.Status {
	return model.NewStatus(status.IsZombie)
}

func (s *statusStore) toEntity(status model.Status) entity.Status {
	return entity.Status{IsZombie: status.Zombie()}
}
