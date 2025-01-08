package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/vandi37/password-manager/internal/application"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	defer stop()

	app := application.New("configs/config.yaml")

	app.Run(ctx)
}
