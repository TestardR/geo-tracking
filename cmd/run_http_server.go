package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/TestardR/geo-tracking/config"
	"github.com/TestardR/geo-tracking/internal/domain/shared"
)

func RunAsHTTPServer(
	cliCtx *cli.Context,
	appVersion string,
	cfg *config.Config,
	consoleOutput shared.StdLogger,
) error {
	ctx := cliCtx.Context

	consoleOutput.Printf("HTTP server mode")

	_, closeRedisStore, err := connectToWithCloseRedisCache(ctx, cfg, consoleOutput)
	if err != nil {
		return err
	}
	defer closeRedisStore()

	// NATS client

	// Http client

	return nil
}
