package application

import (
	"context"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	tgloggerapi "github.com/vandi37/TgLoggerApi"
	"github.com/vandi37/password-manager/internal/config"
	"github.com/vandi37/password-manager/internal/postgresql/database"
	"github.com/vandi37/password-manager/internal/repo/user_repo"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/internal/telegram/commands"
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

	cfg, err := config.Get(a.config)
	if err != nil {
		panic(err)
	}
	out := tgloggerapi.New(cfg.Log.Token, cfg.Log.Chat)
	if !out.FastCheck() {

		panic("can't connect to tg logger")
	}

	logger := logger.New(io.MultiWriter(os.Stderr, out))

	closer := closer.New(logger)

	db, err := database.New(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	if err != nil {
		logger.Fatalln(err)
	}
	closer.Add(db.Close)

	err = db.Init(ctx)
	if err != nil {
		logger.Fatalln(err)
	}

	service := service.New(user_repo.New(db), password.New([]byte(cfg.HashSalt), []byte(cfg.ArgonSalt)))

	b, err := bot.New(cfg.Token, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	b.Init(commands.BuildCommands(b, service, commands.NewUser, commands.UpdateUser, commands.Cancel))

	go b.Run(ctx)

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger.Warnln("Exiting...")
	closer.Close(ctx)
	os.Exit(http.StatusTeapot)
}
