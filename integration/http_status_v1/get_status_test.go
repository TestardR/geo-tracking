package http_status_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bsm/gomega"

	"github.com/TestardR/geo-tracking/config"

	"github.com/TestardR/geo-tracking/integration/components/infrastructure/persistence/redis_cache"
	"github.com/TestardR/geo-tracking/integration/test_shared"
	coordinateService "github.com/TestardR/geo-tracking/internal/application/coordinate_service"
	statusService "github.com/TestardR/geo-tracking/internal/application/status_service"
	httpStatusV1 "github.com/TestardR/geo-tracking/internal/infrastructure/api/http_status_v1"
	"github.com/TestardR/geo-tracking/internal/infrastructure/api/www"
	natsmsEvent "github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms"
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
	coordinateCache "github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/redis_cache"
	"github.com/TestardR/geo-tracking/internal/infrastructure/distance"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
	redisCache "github.com/TestardR/geo-tracking/internal/infrastructure/shared/redis_cache"
	statusCache "github.com/TestardR/geo-tracking/internal/infrastructure/status/redis_cache"
)

func TestGetDriverZombieStatus(t *testing.T) {
	ctx := context.Background()
	redis := redis_cache.MustConnectToRedis(t)
	muteLogger := test_shared.NewMockedSilentLogger(t)
	integrationEnvConfig := test_shared.GetIntegrationConfig(t)

	redisClient := redisCache.NewClient(redis)
	coordinateStore := coordinateCache.NewCoordinateStore(redisClient)

	consumer, err := natsms.NewConsumer(
		integrationEnvConfig.NatsBrokerList,
		"e2e",
		"e2e.test.event",
		muteLogger,
	)
	if err != nil {
		t.Fatal("Error creating consumer")
	}

	distanceFinder := distance.NewDistanceFinder(
		distance.Strategy(distance.HaversineFormula),
		map[distance.Strategy]distance.StrategyExecutor{
			distance.Strategy(distance.HaversineFormula): &distance.Haversine{},
		},
	)
	statusStore := statusCache.NewStatusStore(redisClient)
	statusSvc := statusService.NewService(statusStore, coordinateStore, distanceFinder)

	go consumer.Consume(
		ctx,
		natsmsEvent.NewCoordinateHandler(
			coordinateService.New(coordinateStore),
			statusSvc,
			coordinateStore,
		).Handle,
	)

	producer, err := natsms.NewProducer(
		integrationEnvConfig.NatsBrokerList,
		"e2e",
		"e2e.test.event",
		muteLogger,
	)
	if err != nil {
		t.Fatal("Error creating producer")
	}

	cfg := config.Config(*integrationEnvConfig)
	statusServer := httpStatusV1.NewHttpServer(
		&cfg,
		statusSvc,
		muteLogger,
	)

	ts := httptest.NewServer(statusServer.Handler)

	t.Cleanup(func() {
		ts.Close()
		producer.Close()

		err := redis.Close()
		if err != nil {
			t.Log(err)
		}
	})

	g := gomega.NewWithT(t)

	t.Run("it should return 200 and driver_id status", func(t *testing.T) {
		messages := []entity.DriverCoordinate{
			{
				DriverId:  "123",
				Longitude: 48.908394,
				Latitude:  2.363022,
				CreatedAt: time.Now().UTC(),
			},
			{
				DriverId:  "123",
				Longitude: 48.908261,
				Latitude:  2.364596,
				CreatedAt: time.Now().UTC(),
			},
			{
				DriverId:  "123",
				Longitude: 48.907214,
				Latitude:  2.364462,
				CreatedAt: time.Now().UTC(),
			},
			{
				DriverId:  "123",
				Longitude: 48.907219,
				Latitude:  2.364464,
				CreatedAt: time.Now().UTC(),
			},
		}
		for _, msg := range messages {
			err = producer.Publish(ctx, msg)
			g.Expect(err).To(gomega.BeNil())
		}

		resp, err := http.Get(fmt.Sprintf("%s/status?driver_id=123", ts.URL))
		g.Expect(err).To(gomega.BeNil())

		bodyBytes, err := io.ReadAll(resp.Body)
		g.Expect(err).To(gomega.BeNil())

		var status www.Status
		json.Unmarshal(bodyBytes, &status)
		g.Expect(err).To(gomega.BeNil())

		g.Expect(status.DriverId).To(gomega.Equal("123"))
		g.Expect(status.IsZombie).To(gomega.BeTrue())
		g.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))

	})
}
