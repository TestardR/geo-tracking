package http_status_v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TestardR/geo-tracking/internal/domain/shared"
	httpShared "github.com/TestardR/geo-tracking/internal/infrastructure/http/http-shared"

	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/model"
)

type getStatusHandler interface {
	Handle(ctx context.Context, query query.GetStatus) (model.Status, error)
}

type Response struct {
	DriverId string `json:"driver_id"`
	IsZombie bool   `json:"is_zombie"`
}

type StatusHandler struct {
	getStatusHandler getStatusHandler
}

func NewStatusHandler(getStatusHandler getStatusHandler) *StatusHandler {
	return &StatusHandler{getStatusHandler: getStatusHandler}
}

func (h *StatusHandler) GetStatus(ctx *gin.Context) {
	driverId := ctx.Param("id")
	status, err := h.getStatusHandler.Handle(ctx, query.NewGetStatus(model.NewDriverId(driverId)))
	if err != nil {
		_ = ctx.Error(err)
		if shared.IsDomainError(err) {
			ctx.JSON(http.StatusBadRequest, httpShared.ResponseError{Message: err.Error()})
		}

		ctx.JSON(http.StatusInternalServerError, httpShared.ResponseError{Message: err.Error()})
	}

	ctx.JSON(http.StatusOK, Response{DriverId: driverId, IsZombie: status.Zombie()})
}
