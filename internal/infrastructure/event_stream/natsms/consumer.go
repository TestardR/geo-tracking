package natsms

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

const batchSize = 1

type ConsumerHandler func(context.Context, *nats.Msg) error

type Consumer struct {
	stream     nats.JetStreamContext
	streamName string
	subject    string
	logger     shared.ErrorInfoLogger
	stop       chan struct{}
}

func NewConsumer(
	broker,
	streamName,
	subject string,
	logger shared.ErrorInfoLogger,
) (*Consumer, error) {
	nc, err := nats.Connect(broker)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		stream:     js,
		streamName: streamName,
		subject:    subject,
		logger:     logger,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler ConsumerHandler) error {
	_, err := c.stream.AddStream(&nats.StreamConfig{
		Name:     c.streamName,
		Subjects: []string{c.subject},
	})
	if err != nil {
		return err
	}

	sub, err := c.stream.PullSubscribe("", "", nats.BindStream(c.streamName))
	if err != nil {
		return err
	}

	for {
		select {
		case <-c.stop:
			c.logger.Info(fmt.Sprintf("message consuming for stream %q was stopped", c.streamName))
			err := sub.Unsubscribe()
			if err != nil {
				err = fmt.Errorf(
					"cannot unsubscribe from stream %q: %v",
					c.streamName,
					err,
				)
				c.logger.Error(err)
			}

			return nil
		default:
			messages, err := sub.Fetch(batchSize, nats.MaxWait(time.Second))
			if err != nil {
				if err != nats.ErrTimeout {
					err = fmt.Errorf(
						"cannot read message on stream %q: %v",
						c.streamName,
						err,
					)
					c.logger.Error(err)
				}

				continue
			}

			if len(messages) != batchSize {
				return fmt.Errorf(
					"messages number is greater than %d from stream %q: %v",
					batchSize,
					c.streamName,
					err,
				)
			}

			msg := messages[0]
			err = handler(ctx, msg)
			if err != nil {
				err = fmt.Errorf(
					"cannot handle message from stream %q: %v",
					c.streamName,
					err,
				)
				c.logger.Error(err)

				continue
			}

			err = msg.Ack()
			if err != nil {
				// add message ID ?
				err = fmt.Errorf(
					"cannot acknowledge message from stream %q: %v",
					c.streamName,
					err,
				)
				c.logger.Error(err)
			}

			continue
		}
	}
}
