package natsms

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
)

type coordinateHandler struct {
	coordinatePersister repository.CoordinatePersister
	logger              shared.ErrorLogger
}

func NewCoordinateHandler(
	coordinatePersister repository.CoordinatePersister,
	logger shared.ErrorLogger,
) *coordinateHandler {
	return &coordinateHandler{logger: logger, coordinatePersister: coordinatePersister}
}

func (h *coordinateHandler) Handle(ctx context.Context, msg *nats.Msg) error {
	var driverCoordinate entity.DriverCoordinate
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
