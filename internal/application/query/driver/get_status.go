package driver

import (
	"github.com/TestardR/geo-tracking/internal/domain/driver"
)

type GetStatus struct {
	id     driver.Id
	status driver.Status
}

func NewGetStatus(id driver.Id, status driver.Status) GetStatus {
	return GetStatus{id: id, status: status}
}

func (q GetStatus) Id() driver.Id {
	return q.id
}

func (q GetStatus) Status() driver.Status {
	return q.status
}
