package query

import (
	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type GetStatus struct {
	driverId model.DriverId
	status   model.Status
}

func NewGetStatus(id model.DriverId, status model.Status) GetStatus {
	return GetStatus{driverId: id, status: status}
}

func (q GetStatus) DriverId() model.DriverId {
	return q.driverId
}

func (q GetStatus) Status() model.Status {
	return q.status
}
