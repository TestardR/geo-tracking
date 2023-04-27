package command

import (
	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type AddCoordinate struct {
	driverId   model.DriverId
	coordinate model.Coordinate
}

func NewAddCoordinate(id model.DriverId, coordinate model.Coordinate) AddCoordinate {
	return AddCoordinate{driverId: id, coordinate: coordinate}
}

func (c AddCoordinate) DriverId() model.DriverId {
	return c.driverId
}

func (c AddCoordinate) Coordinate() model.Coordinate {
	return c.coordinate
}
