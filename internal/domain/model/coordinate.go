package model

type Coordinate struct {
	longitude float64
	latitude  float64
}

func NewCoordinate(longitude float64, latitude float64) Coordinate {
	return Coordinate{longitude: longitude, latitude: latitude}
}

func (c *Coordinate) Longitude() float64 {
	return c.longitude
}

func (c *Coordinate) Latitude() float64 {
	return c.latitude
}
