package main

import (
	"context"
	"os/signal"
	"physk/internal/app"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP,
	)
	defer cancel()
	app.App(ctx)
}
