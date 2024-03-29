package http_status_v1

import (
	"context"
	"encoding/json"
	"fmt"
	driverModel "github.com/TestardR/geo-tracking/internal/domain/driver/model"
	"github.com/TestardR/geo-tracking/internal/domain/status/model"
	"net/http"
	"time"

	"github.com/TestardR/geo-tracking/config"
	"github.com/TestardR/geo-tracking/internal/application/query"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/api/www"
)

type getStatusHandler interface {
	HandleGetStatus(ctx context.Context, query query.GetStatus) (model.Status, error)
}

func NewHttpServer(
	cfg *config.Config,
	getStatusHandler getStatusHandler,
	logger shared.ErrorLogger,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.URL.Query().Get("driver_id")
		if id == "" {
			http.Error(w, "request requires a driver_id", http.StatusBadRequest)
			logger.Error("request requires a driver_id as query param")

			return
		}

		driverId := driverModel.NewDriverId(id)
		status, err := getStatusHandler.HandleGetStatus(ctx, query.NewGetStatus(driverId))
		if err != nil {
			if shared.IsDomainError(err) {
				http.Error(w, "driver_id does not exist", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("driver_id does not exist: %v", err))

				return
			}

			http.Error(w, "failed to compute distance", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("failed to compute distance: %v", err))

			return
		}

		data, err := json.Marshal(www.ToWWWStatus(driverId.Id(), status.Zombie()))
		if err != nil {
			http.Error(w, "failed to marshal result", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("failed to marshal result: %v", err))

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	return &http.Server{
		Addr:         cfg.HttpPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
}
