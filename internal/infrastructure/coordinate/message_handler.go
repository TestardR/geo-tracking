package coordinate

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	msgEntity "github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms/entity"
)

type messageHandler struct {
	coordinatePersister repository.CoordinatePersister
	logger              shared.ErrorLogger
}

func NewMessageHandler(
	coordinatePersister repository.CoordinatePersister,
	logger shared.ErrorLogger,
) *messageHandler {
	return &messageHandler{logger: logger, coordinatePersister: coordinatePersister}
}

func (h *messageHandler) Handle(ctx context.Context, msg *nats.Msg) error {
	var driverCoordinate msgEntity.DriverCoordinate
	err := json.Unmarshal(msg.Data, &driverCoordinate)
	if err != nil {
		return err
	}

	driverId := model.NewDriverId(driverCoordinate.DriverId)
	coordinate, err := model.NewCoordinate(driverCoordinate.Longitude, driverCoordinate.Longitude)
	if err != nil {
		return err
	}

	err = h.coordinatePersister.Persist(ctx, driverId, coordinate)
	if err != nil {
		return err
	}

	return nil
}
