package model

import (
	"errors"
)

type Coordinate struct {
	longitude float64
	latitude  float64
}

func NewCoordinate(longitude float64, latitude float64) (Coordinate, error) {
	if longitude < 0 || latitude < 0 {
		return Coordinate{}, errors.New("longitude or latitude cannot be negative")
	}

	return Coordinate{longitude: longitude, latitude: latitude}, nil
}

func (c *Coordinate) Longitude() float64 {
	return c.longitude
}

func (c *Coordinate) Latitude() float64 {
	return c.latitude
}

func RecreateCoordinate(longitude float64, latitude float64) Coordinate {
	return Coordinate{longitude: longitude, latitude: latitude}
}
