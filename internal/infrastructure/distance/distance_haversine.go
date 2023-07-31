package distance

import (
	"context"
	"errors"
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	distanceModel "github.com/TestardR/geo-tracking/internal/domain/distance"
	"github.com/jftuga/geodist"
)

// Haversine computes distance in kilometers between coordinates using Haversine formula
type Haversine struct{}

func (d *Haversine) Distance(ctx context.Context, coordinates []coordinateModel.Coordinate) (distanceModel.Distance, error) {
	if len(coordinates) < 2 {
		return distanceModel.Distance{}, errors.New("coordinates length must be > 1")
	}

	distance := float64(0)
	for i := 0; i < len(coordinates)-1; i++ {
		p1 := geodist.Coord{Lat: coordinates[i].Latitude(), Lon: coordinates[i].Longitude()}
		p2 := geodist.Coord{Lat: coordinates[i+1].Latitude(), Lon: coordinates[i+1].Longitude()}

		_, km := geodist.HaversineDistance(p1, p2)
		distance += km
	}

	return distanceModel.NewDistance(distance), nil
}
