package entity

type DriverCoordinate struct {
	DriverId   string `json:"driver_id"`
	Coordinate struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	}
}
