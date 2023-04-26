package entity

import (
	"time"
)

type DriverCoordinate struct {
	DriverId  string    `json:"driver_id"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
	CreatedAt time.Time `json:"created_at"`
}
