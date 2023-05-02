package natsms

import (
	"context"
	"testing"
	"time"

	"github.com/bsm/gomega"

	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
)

func TestCanProduceMessagesToEventStream(t *testing.T) {
	ctx := context.Background()

	producer := MustConnectToNatsProducer(t, "producer.test", "producer.test.event")

	g := gomega.NewWithT(t)

	t.Run("it can publish to producer", func(t *testing.T) {
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

		producer.Close()
	})
}
