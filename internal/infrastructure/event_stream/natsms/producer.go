package natsms

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms/entity"
)

type Producer struct {
	stream     nats.JetStreamContext
	streamName string
	subject    string
	logger     shared.ErrorInfoLogger
	stop       chan struct{}
}

func NewProducer(
	broker,
	streamName,
	subject string,
	logger shared.ErrorInfoLogger,
) (*Producer, error) {
	nc, err := nats.Connect(broker)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Producer{
		stream:     js,
		streamName: streamName,
		subject:    subject,
		logger:     logger,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, driverCoordinate entity.DriverCoordinate) {
	_, err := p.stream.StreamInfo(p.streamName)
	if err != nats.ErrStreamNotFound {
		p.logger.Error(fmt.Errorf("cannot get nats info: %v", err))

		return
	} else {
		_, err := p.stream.AddStream(&nats.StreamConfig{
			Name:     p.streamName,
			Subjects: []string{p.subject},
		})
		if err != nil {
			p.logger.Error(fmt.Errorf("cannot add nats stream: %v", err))

			return
		}
	}

	for {
		select {
		case <-p.stop:
			p.logger.Info(fmt.Sprintf("message producing to stream %q was stopped", p.streamName))

			return
		default:
			data, err := json.Marshal(driverCoordinate)
			if err != nil {
				p.logger.Error(fmt.Errorf("failed to marshal coordinate: %v", err))

				continue
			}

			_, err = p.stream.Publish(p.subject, data)
			if err != nil {
				p.logger.Error(fmt.Errorf("cannot publish to stream: %v", err))

				continue
			}
		}
	}
}
