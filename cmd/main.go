package main

import (
	"context"
	"flag"
	"github.com/vandi37/password-manager/internal/application"
	"os"
	"os/signal"
)

const (
	config = "config"
)

func main() {
	// Getting config
	cfg := flag.String(config, "configs/config.yaml", "config file path")
	flag.Parse()

	// Notify context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// Running app
	app := application.New(*cfg)
	app.Run(ctx)
}
