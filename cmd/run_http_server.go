package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/TestardR/geo-tracking/config"
	statusService "github.com/TestardR/geo-tracking/internal/application/status_service"
	"github.com/TestardR/geo-tracking/internal/domain/model/distance"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
	httpStatusV1 "github.com/TestardR/geo-tracking/internal/infrastructure/http/http-status-v1"
	logger "github.com/TestardR/geo-tracking/internal/infrastructure/logging/zap_logger"
	redisCache "github.com/TestardR/geo-tracking/internal/infrastructure/persistence/redis_cache"
)

func RunAsHTTPServer(
	cliCtx *cli.Context,
	appVersion string,
	cfg *config.Config,
	appConfig *config.Application,
	consoleOutput shared.StdLogger,
) error {
	ctx := cliCtx.Context

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

	_, err = natsms.NewConsumer(
		cfg.NatsBrokerList,
		natsms.DriverLocationUpdatedStream,
		natsms.DriverLocationUpdatedStream,
		zapSugaredLogger,
	)
	if err != nil {
		return err
	}

	coordinateStore := redisCache.NewCoordinateStore(redisClient)
	distanceFinder := distance.NewDistanceFinder(
		distance.Strategy(appConfig.DistanceAlgorithm),
		map[distance.Strategy]distance.StrategyExecutor{
			"Haversine": &distance.Haversine{},
			"Vincentiy": &distance.Vincenty{},
		},
	)
	statusStore := redisCache.NewStatusStore(
		redisClient,
		coordinateStore,
		distanceFinder,
	)
	statusHandler := httpStatusV1.NewStatusHandler(
		statusService.NewStatus(
			statusStore,
			zapSugaredLogger,
		),
	)
	httpStatusV1.NewStatusHttpServer(cfg, statusHandler)

	return nil
}
