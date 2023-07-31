package command

import (
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	"github.com/TestardR/geo-tracking/internal/domain/driver/model"
)

type ChangeCoordinate struct {
	driverId   model.DriverId
	coordinate coordinateModel.Coordinate
}

func NewChangeCoordinate(id model.DriverId, coordinate coordinateModel.Coordinate) ChangeCoordinate {
	return ChangeCoordinate{driverId: id, coordinate: coordinate}
}

func (c ChangeCoordinate) DriverId() model.DriverId {
	return c.driverId
}

func (c ChangeCoordinate) Coordinate() coordinateModel.Coordinate {
	return c.coordinate
}
