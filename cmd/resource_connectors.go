package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/TestardR/geo-tracking/config"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/system"

	"github.com/redis/go-redis/v9"
)

func connectToWithCloseRedisCache(
	ctx context.Context,
	cfg *config.Config,
	logger shared.StdLogger,
) (*redis.Client, func(), error) {
	var redisClient *redis.Client
	var err error
	connectRetryController := system.NewRetry(10, system.RetryWaitTimeDefault, time.Sleep, logger)

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
		return nil, func() {}, fmt.Errorf("cannot connect to Redis Store: %v", err)
	}

	closeResource := func() {
		logger.Printf("Disconnecting Redis Store")
		err := redisClient.Close()
		if err != nil {
			logger.Printf(fmt.Sprintf("cannot close redis  Micro Store connection: %v", err))
		}
	}

	return redisClient, closeResource, nil
}
