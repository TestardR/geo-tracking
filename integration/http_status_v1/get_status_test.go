package http_status_v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/TestardR/geo-tracking/config"
	"github.com/TestardR/geo-tracking/integration/infrastructure/shared/redis_cache"
	"github.com/TestardR/geo-tracking/integration/test_shared"
	statusService "github.com/TestardR/geo-tracking/internal/application/status_service"
	"github.com/TestardR/geo-tracking/internal/domain/model/distance"
	natsmsEvent "github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms"
	coordinateCache "github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/redis_cache"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
	httpStatusV1 "github.com/TestardR/geo-tracking/internal/infrastructure/http/http_status_v1"
	redisCache "github.com/TestardR/geo-tracking/internal/infrastructure/shared/redis_cache"
	statusCache "github.com/TestardR/geo-tracking/internal/infrastructure/status/redis_cache"
)

func TestGetDriverZombieStatus(t *testing.T) {
	ctx := context.Background()
	redis := redis_cache.MustConnectToRedis(t)
	logger := test_shared.NewMockedLogger(t)
	muteLogger := test_shared.NewMockedSilentLogger(t)
	integrationEnvConfig := test_shared.GetIntegrationConfig(t)

	redisClient := redisCache.NewClient(redis)
	coordinateStore := coordinateCache.NewCoordinateStore(redisClient)

	consumer, err := natsms.NewConsumer(
		integrationEnvConfig.NatsBrokerList,
		natsms.DriverCoordinateUpdatedStream,
		natsms.DriverCoordinateUpdatedSubject,
		muteLogger,
	)
	if err != nil {
		logger.Error("Error creating consumer")
	}

	go consumer.Consume(
		ctx,
		natsmsEvent.NewCoordinateHandler(
			coordinateStore,
			muteLogger,
		).Handle,
	)
	defer consumer.Stop()

	_, err = natsms.NewProducer(
		integrationEnvConfig.NatsBrokerList,
		natsms.DriverCoordinateUpdatedStream,
		natsms.DriverCoordinateUpdatedSubject,
		muteLogger,
	)
	if err != nil {
		logger.Error("Error creating producer")
	}

	distanceFinder := distance.NewDistanceFinder(
		distance.Strategy(distance.HaversineFormula),
		map[distance.Strategy]distance.StrategyExecutor{
			distance.Strategy(distance.HaversineFormula): &distance.Haversine{},
		},
	)
	statusStore := statusCache.NewStatusStore(
		redisClient,
		coordinateStore,
		distanceFinder,
	)

	cfg := config.Config(*integrationEnvConfig)

	statusServer := httpStatusV1.NewStatusHttpServer(
		&cfg,
		statusService.NewStatus(
			statusStore,
			muteLogger,
		),
		muteLogger,
	)
	go func() {
		err := statusServer.ListenAndServe()
		if nil != err && err != http.ErrServerClosed {
			logger.Printf("Server stopped due to the error: %s", err.Error())
		}
	}()

	t.Cleanup(func() {
		if err := statusServer.Shutdown(ctx); err != nil {
			t.Log(err)
		}

		err := redis.Close()
		if err != nil {
			t.Log(err)
		}
	})

}
