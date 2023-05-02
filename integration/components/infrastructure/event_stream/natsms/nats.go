package natsms

import (
	"context"
	"testing"
	"time"

	"github.com/TestardR/geo-tracking/integration/test_shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
	"github.com/TestardR/geo-tracking/internal/infrastructure/system"
)

func MustConnectToNatsProducer(t *testing.T, streamName, subject string) *natsms.Producer {
	connectRetryController := system.NewRetry(15, system.RetryWaitTimeDefault, time.Sleep, test_shared.NewMockedLogger(t))
	var (
		producer *natsms.Producer
		err      error
	)
	err = connectRetryController.Handle(context.Background(), func(ctx context.Context) error {
		producer, err = natsms.NewProducer(
			test_shared.GetIntegrationConfig(t).NatsBrokerList,
			streamName,
			subject,
			test_shared.NewMockedLogger(t),
		)
		return err
	})
	if err != nil {
		panic(err)
	}

	return producer
}
