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
	redisCache "github.com/TestardR/geo-tracking/internal/infrastructure/shared/redis_cache"
	"github.com/TestardR/geo-tracking/internal/infrastructure/status/redis_cache/entity"
)

type statusStore struct {
	redis *redisCache.Client
}

func NewStatusStore(
	redis *redisCache.Client,
) *statusStore {
	return &statusStore{
		redis: redis,
	}
}

func (s *statusStore) Find(ctx context.Context, driverId model.DriverId) (model.Status, error) {
	driverKey := fmt.Sprintf("s:%s", driverId.Id())
	statusAsBytes, err := s.redis.Get(ctx, driverKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Status{}, shared.NewDomainError(repository.ErrDriverIdNotFoundMessage)
		}

		return model.Status{}, err
	}

	if len(statusAsBytes) < 0 {
		return model.Status{}, shared.NewDomainError(repository.ErrDriverIdNotFoundMessage)
	}

	var status entity.Status
	err = json.Unmarshal(statusAsBytes, &status)
	if err != nil {
		return model.Status{}, err
	}

	return model.RecreateStatus(status), nil
}

func (s *statusStore) Persist(ctx context.Context, driverId model.DriverId, status model.Status) error {
	driverKey := fmt.Sprintf("s:%s", driverId.Id())

	statusEntity := statusModelToEntity(status)
	statusAsBytes, err := json.Marshal(statusEntity)
	if err != nil {
		return err
	}

	return s.redis.Set(ctx, driverKey, statusAsBytes, time.Duration(0))
}

func statusModelToEntity(status model.Status) entity.Status {
	return entity.Status{IsZombie: status.Zombie()}
}
