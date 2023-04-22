package www

type Status struct {
	DriverId string `json:"driver_id"`
	IsZombie bool   `json:"is_zombie"`
}

func ToWWWStatus(driverId string, isZombie bool) Status {
	return Status{
		DriverId: driverId,
		IsZombie: isZombie,
	}
}
