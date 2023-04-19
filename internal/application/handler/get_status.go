package handler

import (
	"context"

	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

type Status struct {
	logger shared.ErrorLogger
}

func (s *Status) HandleGetStatus(ctx context.Context, query query.GetStatus) {
	// redis.getStoppedStatus(query.id)
	// status := NewStatus()

	// return status
}
