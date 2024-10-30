package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/mashiike/oteleport"
	"golang.org/x/sys/unix"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, unix.SIGTERM)
	defer stop()

	exitCode, err := oteleport.ClientCLI(ctx, oteleport.ParseClientCLI)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			slog.ErrorContext(ctx, "Failed", "details", err.Error())
		}
	}
	return exitCode
}
