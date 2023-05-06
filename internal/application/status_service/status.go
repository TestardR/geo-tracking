package status_service

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/application/command"
	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
)

type distanceFinder interface {
	Distance(context.Context, []model.Coordinate) (float64, error)
}

type statusFindPersister interface {
	repository.StatusFinder
	repository.StatusPersister
}

type Service struct {
	statusStore     statusFindPersister
	coordinateStore repository.CoordinateFinder
	distanceFinder  distanceFinder
}

func NewService(
	statusStore statusFindPersister,
	coordinateStore repository.CoordinateFinder,
	distanceFinder distanceFinder,
) *Service {
	return &Service{
		statusStore:     statusStore,
		coordinateStore: coordinateStore,
		distanceFinder:  distanceFinder,
	}
}

func (s *Service) HandleGetStatus(ctx context.Context, query query.GetStatus) (model.Status, error) {
	driverId := model.NewDriverId(query.DriverId())
	status, err := s.statusStore.Find(ctx, driverId)
	if err != nil {
		return model.Status{}, err
	}

	return status, nil
}

func (s *Service) HandleChangeStatus(ctx context.Context, cmd command.ChangeStatus) error {
	driverId := model.NewDriverId(cmd.DriverId())
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
