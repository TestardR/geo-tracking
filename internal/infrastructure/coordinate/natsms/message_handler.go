package natsms

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/application/command"
	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
)

type AddCoordinateHandler interface {
	Handle(ctx context.Context, cmd command.AddCoordinate) error
}

type ChangeStatusHandler interface {
	HandleChangeStatus(ctx context.Context, cmd command.ChangeStatus) error
}

type coordinateHandler struct {
	coordinateService AddCoordinateHandler
	statusService     ChangeStatusHandler
	coordinateStore   repository.CoordinateFinder
}

func NewCoordinateHandler(
	coordinateService AddCoordinateHandler,
	statusService ChangeStatusHandler,
	coordinateStore repository.CoordinateFinder,
) *coordinateHandler {
	return &coordinateHandler{
		coordinateService: coordinateService,
		statusService:     statusService,
		coordinateStore:   coordinateStore,
	}
}

func (h *coordinateHandler) Handle(ctx context.Context, msg *nats.Msg) error {
	var driverCoordinate entity.DriverCoordinate
	err := json.Unmarshal(msg.Data, &driverCoordinate)
	if err != nil {
		return err
	}

	driverId := model.NewDriverId(driverCoordinate.DriverId)
	coordinate, err := model.NewCoordinate(
		driverCoordinate.Longitude,
		driverCoordinate.Latitude,
		driverCoordinate.CreatedAt,
	)
	if err != nil {
		return err
	}

	cmd := command.NewAddCoordinate(driverId, coordinate)
	err = h.coordinateService.Handle(ctx, cmd)
	if err != nil {
		return err
	}

	statusCmd := command.NewChangeStatus(driverId)
	err = h.statusService.HandleChangeStatus(ctx, statusCmd)
	if err != nil {
		return err
	}

	return nil
}
