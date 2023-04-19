package driver

import (
	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

type Overview struct {
	id         Id
	coordinate Coordinate
	status     Status
	updatedAt  shared.OccurredAt
}

func NewOverview(
	id Id,
	coordinate Coordinate,
	status Status,
	updatedAt shared.OccurredAt,
) Overview {
	return Overview{
		id:         id,
		coordinate: coordinate,
		status:     status,
		updatedAt:  updatedAt,
	}
}
