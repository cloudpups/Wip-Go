package main

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func main() {
	var configPath string
	var debugMode bool

	app := &cli.App{
		Name:        "wip-go",
		Description: "Do Not Merge, as a service!",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Runs the application",
				Action: func(cCtx *cli.Context) error {
					logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
					zerolog.DefaultContextLogger = &logger

					if debugMode {
						zerolog.SetGlobalLevel(zerolog.DebugLevel)
					}

					config, err := ReadConfig(configPath)

					if err != nil {
						panic(err)
					}

					runApp(*config, logger)

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "config-file",
						Value:       "./config.yaml",
						Usage:       "The path to the configuration file",
						Aliases:     []string{"c"},
						Destination: &configPath,
					},
					&cli.BoolFlag{
						Name:        "debug",
						Destination: &debugMode,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
