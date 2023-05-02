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

type distanceFinder interface {
	Distance(context.Context, []model.Coordinate) (float64, error)
}

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
	distanceFinder    distanceFinder
}

func NewCoordinateHandler(
	coordinateService AddCoordinateHandler,
	statusService ChangeStatusHandler,
	coordinateStore repository.CoordinateFinder,
	distanceFinder distanceFinder,
) *coordinateHandler {
	return &coordinateHandler{
		coordinateService: coordinateService,
		statusService:     statusService,
		coordinateStore:   coordinateStore,
		distanceFinder:    distanceFinder,
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

	coordinates, err := h.coordinateStore.Find(ctx, driverId)
	if err != nil {
		return err
	}

	distance, err := h.distanceFinder.Distance(ctx, coordinates)
	if err != nil {
		return err
	}

	status := model.NewStatus(false)
	status.ComputeZombieStatus(distance)

	statusCmd := command.NewChangeStatus(driverId, status)
	err = h.statusService.HandleChangeStatus(ctx, statusCmd)
	if err != nil {
		return err
	}

	return nil
}
