package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/TestardR/geo-tracking/config"
	coordinateService "github.com/TestardR/geo-tracking/internal/application/coordinate_service"
	statusService "github.com/TestardR/geo-tracking/internal/application/status_service"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	httpStatusV1 "github.com/TestardR/geo-tracking/internal/infrastructure/api/http_status_v1"
	natsmsEvent "github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/natsms"
	coordinateCache "github.com/TestardR/geo-tracking/internal/infrastructure/coordinate/redis_cache"
	"github.com/TestardR/geo-tracking/internal/infrastructure/distance"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
	logger "github.com/TestardR/geo-tracking/internal/infrastructure/logging/zap_logger"
	redisCache "github.com/TestardR/geo-tracking/internal/infrastructure/shared/redis_cache"
	statusCache "github.com/TestardR/geo-tracking/internal/infrastructure/status/redis_cache"
)

func RunAsHTTPServer(
	cliCtx *cli.Context,
	appVersion string,
	cfg *config.Config,
	appConfig *config.Application,
	consoleOutput shared.StdLogger,
) error {
	ctx := cliCtx.Context

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	consoleOutput.Printf("HTTP server mode")

	zapLogger, err := logger.NewLogger(cfg.LogLevel, cfg.LogPath, appVersion, "gt-http-server", cfg.Env)
	if err != nil {
		return err
	}
	zapSugaredLogger := zapLogger.Sugar()

	redisStore, closeRedisStore, err := connectToWithCloseRedisCache(ctx, cfg, consoleOutput)
	if err != nil {
		return err
	}
	defer closeRedisStore()
	redisClient := redisCache.NewClient(redisStore)

	consumer, err := natsms.NewConsumer(
		cfg.NatsBrokerList,
		natsms.DriverCoordinateUpdatedStream,
		natsms.DriverCoordinateUpdatedSubject,
		zapSugaredLogger,
	)
	if err != nil {
		return err
	}

	distanceFinder := distance.NewDistanceFinder(
		distance.Strategy(appConfig.DistanceAlgorithm),
		map[distance.Strategy]distance.StrategyExecutor{
			distance.Strategy(distance.HaversineFormula): &distance.Haversine{},
			distance.Strategy(distance.VincentyFormula):  &distance.Vincenty{},
		},
	)
	coordinateStore := coordinateCache.NewCoordinateStore(redisClient)
	statusStore := statusCache.NewStatusStore(redisClient)
	statusSvc := statusService.NewService(statusStore, coordinateStore, distanceFinder)

	go consumer.Consume(
		ctx,
		natsmsEvent.NewCoordinateHandler(
			coordinateService.New(coordinateStore),
			statusSvc,
			coordinateStore,
		).Handle,
	)
	defer consumer.Stop()

	statusServer := httpStatusV1.NewStatusHttpServer(
		cfg,
		statusSvc,
		zapSugaredLogger,
	)
	go func() {
		consoleOutput.Printf("Starting HTTP Server")
		err := statusServer.ListenAndServe()
		if nil != err && err != http.ErrServerClosed {
			consoleOutput.Printf("Server stopped due to the error: %s", err.Error())
		}
	}()

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		consoleOutput.Printf("Stopping HTTP Server...")
		if err := statusServer.Shutdown(ctx); err != nil {
			consoleOutput.Printf("HTTP Server shutdown error: %s", err.Error())
		}
	}()

	<-stop
	consoleOutput.Printf("Stopping API server")

	return nil
}
