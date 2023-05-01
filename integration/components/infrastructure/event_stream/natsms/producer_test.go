package natsms

import (
	"context"
	"testing"
	"time"

	"github.com/bsm/gomega"

	"github.com/TestardR/geo-tracking/integration/test_shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
)

func TestCanProduceMessagesToEventStream(t *testing.T) {
	ctx := context.Background()
	logger := test_shared.NewMockedLogger(t)

	muteLogger := test_shared.NewMockedSilentLogger(t)
	integrationEnvConfig := test_shared.GetIntegrationConfig(t)

	producer, err := natsms.NewProducer(
		integrationEnvConfig.NatsBrokerList,
		natsms.DriverCoordinateUpdatedStream,
		natsms.DriverCoordinateUpdatedSubject,
		muteLogger,
	)
	if err != nil {
		logger.Error("Error creating producer")
	}

	g := gomega.NewWithT(t)

	t.Run("it can publish to producer", func(t *testing.T) {
		messages := []entity.DriverCoordinate{
			{
				DriverId:  "123",
				Longitude: 48.908394,
				Latitude:  2.363022,
				CreatedAt: time.Now().UTC(),
			},
		}

		for _, msg := range messages {
			err := producer.Publish(ctx, msg)
			g.Expect(err).To(gomega.BeNil())
		}

	})
}
