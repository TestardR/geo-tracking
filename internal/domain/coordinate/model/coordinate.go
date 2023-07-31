package model

import (
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
) Coordinate {
	return Coordinate{
		longitude: longitude,
		latitude:  latitude,
		createdAt: createdAt,
	}
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
