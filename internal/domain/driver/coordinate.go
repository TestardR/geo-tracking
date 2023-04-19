package driver

type Coordinate struct {
	longitude int64
	latitude  int64
}

func NewCoordinate(longitude int64, latitude int64) Coordinate {
	return Coordinate{longitude: longitude, latitude: latitude}
}
