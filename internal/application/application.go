package application

import (
	"context"
	"fmt"
	"github.com/vandi37/password-manager/internal/repo/password_repo"
	"github.com/vandi37/password-manager/internal/repo/user_repo"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/vandi37/password-manager/internal/config"
	"github.com/vandi37/password-manager/internal/postgresql/database"
	"github.com/vandi37/password-manager/internal/repo/repo"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/internal/telegram/commands"
	"github.com/vandi37/password-manager/internal/telegram/commands/password_commands"
	"github.com/vandi37/password-manager/internal/telegram/commands/user_commands"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/closer"
	"github.com/vandi37/password-manager/pkg/logger"
	"github.com/vandi37/password-manager/pkg/password"
)

type Application struct {
	config string
	mu     sync.Mutex
}

func New(config string) *Application {
	return &Application{
		config: config,
	}
}

func (a *Application) Run(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	l := logger.ConsoleAndFile("logs.log")

	cfg, err := config.Get(a.config)
	if err != nil {
		l.Fatal("failed to load config file", zap.Error(err))
	}

	cl := closer.New()

	db, err := database.New(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	if err != nil {
		l.Fatal("failed to connect database", zap.Error(err))
	}
	cl.Add(db.Close)

	err = db.Init(ctx)
	if err != nil {
		l.Fatal("failed to connect database", zap.Error(err))
	}

	l.Info("database connection established")

	s := service.New(repo.New(user_repo.New(db), password_repo.New(db)), password.New([]byte(cfg.HashSalt), []byte(cfg.ArgonSalt)))

	b, err := bot.New(cfg.Token)
	if err != nil {
		l.Info("failed to create telegram bot", zap.Error(err))
	}

	b.Init(commands.BuildCommands(b, s, user_commands.NewUser, user_commands.UpdateUser, user_commands.DeleteUser, password_commands.ViewByUser, password_commands.NewPassword, password_commands.ViewByCompany, commands.Cancel, commands.GeneratePassword, commands.Help, commands.Start))

	go b.Run(logger.Context(ctx, l))
	l.Info(fmt.Sprintf("bot started at %s", b.GetUsername()))

	<-ctx.Done()
	l.Info("Shutting down service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = cl.Close(ctx)
	if err != nil {
		l.Fatal("failed to close all", zap.Error(err))
	}
}
