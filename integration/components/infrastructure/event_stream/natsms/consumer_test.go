package natsms

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/bsm/gomega"
	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/integration/test_shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
)

func TestCanConsumeMessagesFromEventStream(t *testing.T) {
	ctx := context.Background()
	muteLogger := test_shared.NewMockedSilentLogger(t)
	integrationEnvConfig := test_shared.GetIntegrationConfig(t)

	producer := MustConnectToNatsProducer(t, "consumer", "consumer.test.event")
	consumer, err := natsms.NewConsumer(
		integrationEnvConfig.NatsBrokerList,
		"consumer",
		"consumer.test.event",
		muteLogger,
	)
	if err != nil {
		t.Fatal("Error creating consumer", err)
	}

	g := gomega.NewWithT(t)

	t.Run("it can consume from event stream", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			msg := entity.DriverCoordinate{
				DriverId:  "123",
				Longitude: 48.908394,
				Latitude:  2.363022,
				CreatedAt: time.Now().UTC(),
			}

			err := producer.Publish(ctx, msg)
			g.Expect(err).To(gomega.BeNil())
		}

		consumeCh := make(chan *nats.Msg)
		go func() {
			consumer.Consume(context.Background(), func(ctx context.Context, msg *nats.Msg) error {
				consumeCh <- msg

				return nil
			})
		}()

		counter := 0
		for msg := range consumeCh {
			counter++

			var event entity.DriverCoordinate
			err := json.Unmarshal(msg.Data, &event)
			g.Expect(err).To(gomega.BeNil())

			g.Expect(event.DriverId).To(gomega.Equal("123"))
			g.Expect(event.Longitude).To(gomega.Equal(48.908394))
			g.Expect(event.Latitude).To(gomega.Equal(2.363022))

			if counter == 3 {
				break
			}
		}

		producer.Close()
	})
}
