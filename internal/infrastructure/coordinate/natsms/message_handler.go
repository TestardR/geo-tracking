package natsms

import (
	"context"
	"encoding/json"
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	coordinateRepository "github.com/TestardR/geo-tracking/internal/domain/coordinate/repository"
	"github.com/TestardR/geo-tracking/internal/domain/driver/model"
	"github.com/nats-io/nats.go"

	"github.com/TestardR/geo-tracking/internal/application/command"
	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms/entity"
)

type AddCoordinateHandler interface {
	Handle(ctx context.Context, cmd command.ChangeCoordinate) error
}

type ChangeStatusHandler interface {
	HandleChangeStatus(ctx context.Context, cmd command.ChangeStatus) error
}

type coordinateHandler struct {
	coordinateService AddCoordinateHandler
	statusService     ChangeStatusHandler
	coordinateStore   coordinateRepository.CoordinateFinder
}

func NewCoordinateHandler(
	coordinateService AddCoordinateHandler,
	statusService ChangeStatusHandler,
	coordinateStore coordinateRepository.CoordinateFinder,
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
	coordinate := coordinateModel.NewCoordinate(
		driverCoordinate.Longitude,
		driverCoordinate.Latitude,
		driverCoordinate.CreatedAt,
	)

	cmd := command.NewChangeCoordinate(driverId, coordinate)
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
