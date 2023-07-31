package coordinate_service

import (
	"context"
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	coordinateRepository "github.com/TestardR/geo-tracking/internal/domain/coordinate/repository"
	"github.com/TestardR/geo-tracking/internal/domain/coordinate/validator"

	"github.com/TestardR/geo-tracking/internal/application/command"
)

type Service struct {
	coordinatePersister coordinateRepository.CoordinatePersister
	coordinateValidator validator.CoordinateValidator
}

func New(coordinatePersister coordinateRepository.CoordinatePersister) *Service {
	return &Service{coordinatePersister: coordinatePersister}
}

func (s *Service) Handle(ctx context.Context, cmd command.ChangeCoordinate) error {
	addCoordinateChange, err := coordinateModel.AddCoordinate(
		ctx,
		cmd.DriverId(),
		cmd.Coordinate(),
		s.coordinateValidator,
	)
	if err != nil {
		return err
	}

	err = s.coordinatePersister.Persist(ctx, addCoordinateChange)
	if err != nil {
		return err
	}

	return nil
}
