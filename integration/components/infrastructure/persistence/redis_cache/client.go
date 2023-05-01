package redis_cache

import (
	"context"
	"testing"
	"time"

	"github.com/TestardR/geo-tracking/integration/test_shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/system"
	"github.com/redis/go-redis/v9"
)

func MustConnectToRedis(t *testing.T) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cfg := test_shared.GetIntegrationConfig(t)
	connectRetryController := system.NewRetry(15, system.RetryWaitTimeDefault, time.Sleep, test_shared.NewMockedLogger(t))

	var (
		redisClient *redis.Client
		err         error
	)

	err = connectRetryController.Handle(ctx, func(ctx context.Context) error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: cfg.RedisMasterAddr,
			DB:   cfg.RedisDb,
		})
		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return redisClient
}
