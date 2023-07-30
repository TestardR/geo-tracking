package model

type DriverId struct {
	id string
}

func NewDriverId(id string) DriverId {
	return DriverId{id: id}
}

func (i DriverId) Id() string {
	return i.id
}
