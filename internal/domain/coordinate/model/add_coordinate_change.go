package model

import (
	"context"
	"github.com/TestardR/geo-tracking/internal/domain/coordinate/validator"
	"github.com/TestardR/geo-tracking/internal/domain/driver/model"
)

type AddCoordinateChange struct {
	driverId   model.DriverId
	coordinate Coordinate
}

func NewAddCoordinateChange(
	driverId model.DriverId,
	coordinate Coordinate,
) AddCoordinateChange {
	return AddCoordinateChange{
		driverId:   driverId,
		coordinate: coordinate,
	}
}

func (c AddCoordinateChange) Coordinate() Coordinate {
	return c.coordinate
}

func (c AddCoordinateChange) DriverId() model.DriverId {
	return c.driverId
}

func AddCoordinate(
	ctx context.Context,
	driverId model.DriverId,
	coordinate Coordinate,
	validator validator.CoordinateValidator,
) (AddCoordinateChange, error) {
	err := validator.CoordinateValid(ctx, coordinate.Longitude(), coordinate.Latitude())
	if err != nil {
		return AddCoordinateChange{}, err
	}

	return NewAddCoordinateChange(driverId, coordinate), nil
}
