package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/TestardR/geo-tracking/config"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
	"github.com/TestardR/geo-tracking/internal/infrastructure/event_stream/natsms"
	"github.com/TestardR/geo-tracking/internal/infrastructure/logging/zap_logger"
)

func RunAsHTTPServer(
	cliCtx *cli.Context,
	appVersion string,
	cfg *config.Config,
	consoleOutput shared.StdLogger,
) error {
	ctx := cliCtx.Context

	consoleOutput.Printf("HTTP server mode")

	zapLogger, err := zap_logger.NewLogger(cfg.LogLevel, cfg.LogPath, appVersion, "gt-http-server", cfg.Env)
	if err != nil {
		return err
	}
	zapSugaredLogger := zapLogger.Sugar()

	_, closeRedisStore, err := connectToWithCloseRedisCache(ctx, cfg, consoleOutput)
	if err != nil {
		return err
	}
	defer closeRedisStore()

	_, err = natsms.NewConsumer(
		cfg.NatsBrokerList,
		natsms.GeoLocationStream,
		natsms.GeoLocationStream,
		zapSugaredLogger,
	)
	if err != nil {
		return err
	}

	// Http client

	return nil
}
