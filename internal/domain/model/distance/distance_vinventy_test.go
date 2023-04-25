package distance

import (
	"context"
	"testing"

	"github.com/bsm/gomega"

	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/infrastructure/persistence/redis_cache/entity"
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

		var coordinatesModel []model.Coordinate
		for _, coordinate := range coordinates {
			coordinatesModel = append(coordinatesModel, model.RecreateCoordinate(coordinate.Longitude, coordinate.Latitude))
		}

		result, _ := distance.Distance(context.Background(), coordinatesModel)
		g.Expect(result).To(gomega.Equal(0.11739179888295599))
	})

	t.Run("can compute coordinates with even values", func(t *testing.T) {
		coordinates := []entity.Coordinate{
			{Longitude: 48.908394, Latitude: 2.363022},
			{Longitude: 48.908261, Latitude: 2.364596},
			{Longitude: 48.907214, Latitude: 2.364462},
			{Longitude: 48.906223, Latitude: 2.364355},
		}

		var coordinatesModel []model.Coordinate
		for _, coordinate := range coordinates {
			coordinatesModel = append(coordinatesModel, model.RecreateCoordinate(coordinate.Longitude, coordinate.Latitude))
		}

		result, _ := distance.Distance(context.Background(), coordinatesModel)
		g.Expect(result).To(gomega.Equal(0.22824931425778677))
	})
}
