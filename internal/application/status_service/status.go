package status_service

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/persistence/redis_cache/entity"
)

type statusFinder interface {
	Find(ctx context.Context, driver entity.Driver) (model.Status, error)
}

type Status struct {
	statusFinder statusFinder
	logger       shared.ErrorLogger
}

func NewStatus(
	statusFinder statusFinder,
	logger shared.ErrorLogger,
) *Status {
	return &Status{logger: logger, statusFinder: statusFinder}
}

func (s *Status) Handle(ctx context.Context, query query.GetStatus) (model.Status, error) {
	driver := entity.Driver{
		Id: query.DriverId(),
	}
	status, err := s.statusFinder.Find(ctx, driver)
	if err != nil {
		return model.Status{}, err
	}

	return status, nil
}
