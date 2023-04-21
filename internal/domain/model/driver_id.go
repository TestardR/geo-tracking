package model

type DriverId struct {
	Id string
}

func NewDriverId(id string) DriverId {
	return DriverId{Id: id}
}
