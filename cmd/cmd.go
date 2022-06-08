package cmd

import (
	"context"
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"
)

func Execute(ctx context.Context, cancel context.CancelFunc, args []string) {
	app := &cli.App{
		Name:  "iNotify Automation",
		Usage: "to make system actions on FS events",
		Commands: []*cli.Command{
			start(),
		},
	}

	err := app.RunContext(ctx, args)
	if err != nil {
		fmt.Printf("execution failed: %s", err)
		cancel()
		os.Exit(1)
	}

	cancel()
}
