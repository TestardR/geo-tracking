package natsms

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

type ConsumerHandler func(m *nats.Msg)

type Consumer struct {
	topic  string
	logger shared.ErrorInfoLogger
}

func NewConsumer(
	broker,
	topic string,
	subject string,
	logger shared.ErrorInfoLogger,
) (*Consumer, error) {
	//nc, err := nats.Connect(broker)
	//if err != nil {
	//	return nil, err
	//}

	return &Consumer{
		topic:  topic,
		logger: logger,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler ConsumerHandler) {

}
