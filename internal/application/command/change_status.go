package command

import (
	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type ChangeStatus struct {
	driverId model.DriverId
}

func NewChangeStatus(id model.DriverId) ChangeStatus {
	return ChangeStatus{driverId: id}
}

func (c ChangeStatus) DriverId() string {
	return c.driverId.Id()
}
