package status_service

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

type Status struct {
	statusFinder repository.StatusFinder
	logger       shared.ErrorLogger
}

func NewStatus(
	statusFinder repository.StatusFinder,
	logger shared.ErrorLogger,
) *Status {
	return &Status{logger: logger, statusFinder: statusFinder}
}

func (s *Status) Handle(ctx context.Context, query query.GetStatus) (model.Status, error) {
	driverId := model.NewDriverId(query.DriverId())
	status, err := s.statusFinder.Find(ctx, driverId)
	if err != nil {
		return model.Status{}, err
	}

	return status, nil
}
