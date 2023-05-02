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

func TestCanConsumeMessagesFromEventStream(t *testing.T) {
	ctx := context.Background()
	muteLogger := test_shared.NewMockedSilentLogger(t)
	integrationEnvConfig := test_shared.GetIntegrationConfig(t)

	producer, err := natsms.NewProducer(
		integrationEnvConfig.NatsBrokerList,
		"consumer.test",
		"consumer.test.event",
		muteLogger,
	)
	if err != nil {
		t.Fatal("Error creating producer")
	}

	consumer, err := natsms.NewConsumer(
		integrationEnvConfig.NatsBrokerList,
		"consumer.test",
		"consumer.test.event",
		muteLogger,
	)
	if err != nil {
		t.Fatal("Error creating consumer")
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
