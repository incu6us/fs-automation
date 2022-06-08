package cmd

import (
	"os"
	"path/filepath"

	cli "github.com/urfave/cli/v2"

	"github.com/incu6us/fs-automation/executor"
)

const (
	configFlagName = "cfg"
)

const (
	configFileName = "config.yaml"
)

func start() *cli.Command {
	command := &cli.Command{
		Name:        "start",
		Description: "run the application",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:      configFlagName,
				Usage:     "path to configuration. by Default it will take the file from app directory",
				Aliases:   []string{"c"},
				TakesFile: true,
			},
		},
		Action: func(c *cli.Context) error {
			filePath, err := configFilePath(c)
			if err != nil {
				return err
			}

			cfg, err := executor.NewConfig(filePath)
			if err != nil {
				return err
			}

			svc := executor.NewService(cfg)
			err = svc.Run(c.Context)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return command
}

func configFilePath(c *cli.Context) (string, error) {
	if c.IsSet(configFlagName) {
		return c.Path(configFlagName), nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, configFileName), nil
}
