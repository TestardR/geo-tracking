package command

import (
	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type ChangeStatus struct {
	driverId model.DriverId
	status   model.Status
}

func NewChangeStatus(id model.DriverId, status model.Status) ChangeStatus {
	return ChangeStatus{driverId: id, status: status}
}

func (c ChangeStatus) DriverId() string {
	return c.driverId.Id()
}

func (c ChangeStatus) Status() model.Status {
	return c.status
}
