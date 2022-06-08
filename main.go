package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/incu6us/fs-automation/cmd"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func(cancelFunc context.CancelFunc) {
		select {
		case sign := <-signalChan:
			fmt.Printf("got %s signal", sign)
			cancel()
		case <-ctx.Done():
			return
		}
	}(cancel)

	cmd.Execute(ctx, cancel, os.Args)
}
