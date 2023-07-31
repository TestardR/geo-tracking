package query

import (
	"github.com/TestardR/geo-tracking/internal/domain/driver/model"
)

type GetStatus struct {
	driverId model.DriverId
}

func NewGetStatus(id model.DriverId) GetStatus {
	return GetStatus{driverId: id}
}

func (q GetStatus) DriverId() string {
	return q.driverId.Id()
}
