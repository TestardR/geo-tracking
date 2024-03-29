package distance

import (
	"context"
	coordinateModel "github.com/TestardR/geo-tracking/internal/domain/coordinate/model"
	"testing"
	"time"

	"github.com/bsm/gomega"

	"github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/redis_cache/entity"
)

func TestCanComputeDistanceWithVincentyFormula(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)
	distance := Vincenty{}

	t.Run("can compute coordinates with odd values", func(t *testing.T) {
		coordinates := []entity.Coordinate{
			{Longitude: 48.908394, Latitude: 2.363022},
			{Longitude: 48.908261, Latitude: 2.364596},
			{Longitude: 48.907214, Latitude: 2.364462},
		}

		var coordinatesModel []coordinateModel.Coordinate
		for _, coordinate := range coordinates {
			coordinatesModel = append(coordinatesModel, coordinateModel.RecreateCoordinate(coordinate.Longitude, coordinate.Latitude, time.Now()))
		}

		result, _ := distance.Distance(context.Background(), coordinatesModel)
		g.Expect(result.Kilometers()).To(gomega.Equal(0.11739179888295599))
	})

	t.Run("can compute coordinates with even values", func(t *testing.T) {
		coordinates := []entity.Coordinate{
			{Longitude: 48.908394, Latitude: 2.363022},
			{Longitude: 48.908261, Latitude: 2.364596},
			{Longitude: 48.907214, Latitude: 2.364462},
			{Longitude: 48.906223, Latitude: 2.364355},
		}

		var coordinatesModel []coordinateModel.Coordinate
		for _, coordinate := range coordinates {
			coordinatesModel = append(coordinatesModel, coordinateModel.RecreateCoordinate(coordinate.Longitude, coordinate.Latitude, time.Now()))
		}

		result, _ := distance.Distance(context.Background(), coordinatesModel)
		g.Expect(result.Kilometers()).To(gomega.Equal(0.22824931425778677))
	})
}
