package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/vandi37/password-manager/internal/application"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGABRT, syscall.SIGALRM, syscall.SIGBUS, syscall.SIGFPE, syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT, syscall.SIGPIPE, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGTERM, syscall.SIGTRAP)
	defer stop()

	app := application.New("configs/config.yaml")

	app.Run(ctx)
}
