package distance

import (
	"context"
	"errors"

	"github.com/jftuga/geodist"

	"github.com/TestardR/geo-tracking/internal/domain/model"
)

// Distance computes distance in kilometers between coordinates using Haversine formula
type distanceHaversine struct{}

func (d *distanceHaversine) Distance(ctx context.Context, coordinates []model.Coordinate) (float64, error) {
	if len(coordinates) < 2 {
		return 0, errors.New("coordinates length must be > 1")
	}

	distance := float64(0)
	for i := 1; i <= len(coordinates)-1; i++ {
		p1 := geodist.Coord{Lat: coordinates[i].Latitude(), Lon: coordinates[i].Longitude()}
		p2 := geodist.Coord{Lat: coordinates[i+1].Latitude(), Lon: coordinates[i+1].Longitude()}

		_, km := geodist.HaversineDistance(p1, p2)
		distance += km
	}

	return distance, nil
}
