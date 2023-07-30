package validator

import (
	"context"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

type CoordinateValidator struct{}

func NewCoordinateValidator() CoordinateValidator {
	return CoordinateValidator{}
}

func (v CoordinateValidator) CoordinateValid(ctx context.Context, longitude, latitude float64) error {
	if longitude < 0 || latitude < 0 {
		return shared.NewDomainError("longitude or latitude cannot be negative")
	}

	return nil
}
