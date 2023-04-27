package coordinate_service

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/application/command"
	"github.com/TestardR/geo-tracking/internal/domain/repository"
)

type Service struct {
	coordinatePersister repository.CoordinatePersister
}

func New(coordinatePersister repository.CoordinatePersister) *Service {
	return &Service{coordinatePersister: coordinatePersister}
}

func (s *Service) Handle(ctx context.Context, cmd command.AddCoordinate) error {
	err := s.coordinatePersister.Persist(ctx, cmd.DriverId(), cmd.Coordinate())
	if err != nil {
		return err
	}

	return nil
}
