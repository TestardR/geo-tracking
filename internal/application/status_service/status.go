package status_service

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
)

type Service struct {
	statusFinder repository.StatusFinder
}

func New(statusFinder repository.StatusFinder) *Service {
	return &Service{statusFinder: statusFinder}
}

func (s *Service) Handle(ctx context.Context, query query.GetStatus) (model.Status, error) {
	driverId := model.NewDriverId(query.DriverId())
	status, err := s.statusFinder.Find(ctx, driverId)
	if err != nil {
		return model.Status{}, err
	}

	return status, nil
}
