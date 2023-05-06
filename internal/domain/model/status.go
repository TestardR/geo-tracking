package model

import (
	"github.com/TestardR/geo-tracking/internal/infrastructure/status/redis_cache/entity"
)

type Status struct {
	isZombie bool
}

func NewStatus() Status {
	return Status{}
}

func (s *Status) Zombie() bool {
	return s.isZombie
}

func NewStatusFromDistance(distance float64) Status {
	isZombie := false
	if distance <= 0.5 {
		isZombie = true
	}

	return Status{isZombie: isZombie}
}

func RecreateStatus(entity entity.Status) Status {
	return Status{
		isZombie: entity.IsZombie,
	}
}
