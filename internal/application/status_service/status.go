package status_service

import (
	"context"
	"github.com/TestardR/geo-tracking/internal/application/command"
	"github.com/TestardR/geo-tracking/internal/application/query"
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	coordinateRepository "github.com/TestardR/geo-tracking/internal/domain/coordinate/repository"
	distanceModel "github.com/TestardR/geo-tracking/internal/domain/distance"
	driverModel "github.com/TestardR/geo-tracking/internal/domain/driver/model"
	"github.com/TestardR/geo-tracking/internal/domain/status/model"
	"github.com/TestardR/geo-tracking/internal/domain/status/repository"
)

type distanceFinder interface {
	Distance(context.Context, []coordinateModel.Coordinate) (distanceModel.Distance, error)
}

type statusFindPersister interface {
	repository.StatusFinder
	repository.StatusPersister
}

type Service struct {
	statusStore     statusFindPersister
	coordinateStore coordinateRepository.CoordinateFinder
	distanceFinder  distanceFinder
}

func NewService(
	statusStore statusFindPersister,
	coordinateStore coordinateRepository.CoordinateFinder,
	distanceFinder distanceFinder,
) *Service {
	return &Service{
		statusStore:     statusStore,
		coordinateStore: coordinateStore,
		distanceFinder:  distanceFinder,
	}
}

func (s *Service) HandleGetStatus(ctx context.Context, query query.GetStatus) (model.Status, error) {
	driverId := driverModel.NewDriverId(query.DriverId())
	status, err := s.statusStore.Find(ctx, driverId)
	if err != nil {
		return model.Status{}, err
	}

	return status, nil
}

func (s *Service) HandleChangeStatus(ctx context.Context, cmd command.ChangeStatus) error {
	driverId := driverModel.NewDriverId(cmd.DriverId())
	coordinates, err := s.coordinateStore.Find(ctx, driverId)
	if err != nil {
		return err
	}

	distance, err := s.distanceFinder.Distance(ctx, coordinates)
	if err != nil {
		return err
	}

	status := model.NewStatusFromDistance(distance)

	err = s.statusStore.Persist(ctx, driverId, status)
	if err != nil {
		return err
	}

	return nil
}
