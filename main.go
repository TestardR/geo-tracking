package main

import (
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/TestardR/geo-tracking/cmd"
	"github.com/TestardR/geo-tracking/config"
)

var version = "develop"

func main() {
	consoleOutput := log.New(os.Stdout, "", log.LstdFlags)
	consoleOutput.Printf("Available processors: %d", runtime.NumCPU())
	consoleOutput.Printf("GoMaxProcs: %d", runtime.GOMAXPROCS(0))

	cfg, err := config.NewConfig()
	if err != nil {
		consoleOutput.Fatal(err)
	}

	appCfg, err := config.NewApplication()
	if err != nil {
		consoleOutput.Fatal(err)
	}

	app := &cli.App{
		Name:  "geo-tracking",
		Usage: "CLI application for managing GT services",
		Commands: []*cli.Command{
			{
				Name:    "http-server",
				Usage:   "Starts driver status HTTP server",
				Aliases: []string{"s"},
				Action: func(c *cli.Context) error {
					return cmd.RunAsHTTPServer(c, version, cfg, appCfg, consoleOutput)
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		consoleOutput.Fatal(err)
	}
}
