package model

import (
	"github.com/TestardR/geo-tracking/internal/infrastructure/status/redis_cache/entity"
)

type Status struct {
	isZombie bool
}

func NewStatus(isZombie bool) Status {
	return Status{isZombie: isZombie}
}

func (s *Status) Zombie() bool {
	return s.isZombie
}

func (s *Status) ComputeZombieStatus(distance float64) bool {
	if distance <= 0.5 {
		s.isZombie = true
	}

	return s.isZombie
}

func StatusToEntity(status Status) entity.Status {
	return entity.Status{
		IsZombie: status.Zombie(),
	}
}
