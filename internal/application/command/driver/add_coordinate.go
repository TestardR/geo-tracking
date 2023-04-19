package driver

import (
	"github.com/TestardR/geo-tracking/internal/domain/driver"
)

type AddCoordinate struct {
	id         driver.Id
	coordinate driver.Coordinate
}

func NewAddCoordinate(id driver.Id, coordinate driver.Coordinate) AddCoordinate {
	return AddCoordinate{id: id, coordinate: coordinate}
}

func (c AddCoordinate) Id() driver.Id {
	return c.id
}

func (c AddCoordinate) Name() driver.Coordinate {
	return c.coordinate
}
