package model

import (
	"errors"
	"time"
)

type Coordinate struct {
	longitude float64
	latitude  float64
	createdAt time.Time
}

func NewCoordinate(
	longitude float64,
	latitude float64,
	createdAt time.Time,
) (Coordinate, error) {
	if longitude < 0 || latitude < 0 {
		return Coordinate{}, errors.New("longitude or latitude cannot be negative")
	}

	return Coordinate{
		longitude: longitude,
		latitude:  latitude,
		createdAt: createdAt,
	}, nil
}

func (c *Coordinate) Longitude() float64 {
	return c.longitude
}

func (c *Coordinate) Latitude() float64 {
	return c.latitude
}

func (c *Coordinate) CreatedAt() time.Time {
	return c.createdAt
}

func RecreateCoordinate(
	longitude float64,
	latitude float64,
	createdAt time.Time,
) Coordinate {
	return Coordinate{
		longitude: longitude,
		latitude:  latitude,
		createdAt: createdAt,
	}
}
