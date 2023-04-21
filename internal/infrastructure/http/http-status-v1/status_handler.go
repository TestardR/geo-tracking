package http_status_v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	sharedError "github.com/TestardR/geo-tracking/internal/domain/shared"

	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/model"
	"github.com/TestardR/geo-tracking/internal/infrastructure/http/shared"
)

type getStatusHandler interface {
	HandleGetStatus(ctx context.Context, query query.GetStatus) (model.Status, error)
}

type StatusHandler struct {
	getStatusHandler getStatusHandler
}

type Response struct {
	DriverId string `json:"driver_id"`
	IsZombie bool   `json:"is_zombie"`
}

func (h *StatusHandler) GetStatus(ctx *gin.Context) {
	driverId := ctx.Param("id")
	status, err := h.getStatusHandler.HandleGetStatus(ctx, query.NewGetStatus(model.NewDriverId(driverId)))
	if err != nil {
		_ = ctx.Error(err)
		if sharedError.IsDomainError(err) {
			ctx.JSON(http.StatusBadRequest, shared.ResponseError{Message: err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, shared.ResponseError{Message: err.Error()})
	}

	ctx.JSON(http.StatusOK, Response{DriverId: driverId, IsZombie: status.Zombie()})

}
