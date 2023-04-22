package coordinate_worker

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/domain/shared"
	msgEntity "github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms/entity"
	"github.com/TestardR/geo-tracking/internal/infrastructure/persistence/redis_cache/entity"
)

type coordinatePersister interface {
	Persist(ctx context.Context, driver entity.Driver, coordinate entity.Coordinate) error
}

type messageHandler struct {
	coordinatePersister coordinatePersister
	logger              shared.ErrorLogger
}

func NewMessageHandler(
	coordinatePersister coordinatePersister,
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

	driver := entity.Driver{Id: driverCoordinate.DriverId}
	coordinate := entity.Coordinate{
		Longitude: driverCoordinate.Coordinate.Longitude,
		Latitude:  driverCoordinate.Coordinate.Latitude,
	}
	err = h.coordinatePersister.Persist(ctx, driver, coordinate)
	if err != nil {
		return err
	}

	return nil
}
