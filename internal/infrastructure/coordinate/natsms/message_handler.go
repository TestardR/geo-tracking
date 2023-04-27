package natsms

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/application/command"
	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
)

type Handler interface {
	Handle(ctx context.Context, cmd command.AddCoordinate) error
}

type coordinateHandler struct {
	coordinateService Handler
}

func NewCoordinateHandler(coordinateService Handler) *coordinateHandler {
	return &coordinateHandler{coordinateService: coordinateService}
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

	cmd := command.NewAddCoordinate(driverId, coordinate)
	err = h.coordinateService.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
