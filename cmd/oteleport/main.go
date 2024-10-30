package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/mashiike/oteleport"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	var sigterm os.Signal = syscall.SIGTERM
	if runtime.GOOS == "windows" {
		sigterm = os.Interrupt
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, sigterm)
	defer stop()

	exitCode, err := oteleport.ServerCLI(ctx, oteleport.ParseServerCLI)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			slog.ErrorContext(ctx, "Failed", "details", err.Error())
		}
	}
	return exitCode
}
