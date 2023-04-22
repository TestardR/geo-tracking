package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/TestardR/geo-tracking/config"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms/entity"
	logger "github.com/TestardR/geo-tracking/internal/infrastructure/logging/zap_logger"
)

func RunCoordinateWorker(
	cliCtx *cli.Context,
	appVersion string,
	cfg *config.Config,
	appConfig *config.Application,
	consoleOutput shared.StdLogger,
) error {
	ctx := cliCtx.Context

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	consoleOutput.Printf("Coordinate worker started")

	zapLogger, err := logger.NewLogger(cfg.LogLevel, cfg.LogPath, appVersion, "coordiante-worker", cfg.Env)
	if err != nil {
		return err
	}
	zapSugaredLogger := zapLogger.Sugar()

	producer, err := natsms.NewProducer(
		cfg.NatsBrokerList,
		natsms.DriverCoordindateUpdatedStream,
		natsms.DriverCoordindateUpdatedStream,
		zapSugaredLogger,
	)
	if err != nil {
		return err
	}

	// TODO: produce entity
	producer.Publish(ctx, entity.DriverCoordinate{})

	return nil
}
