package bot

import (
	"context"
	"github.com/vandi37/password-manager/pkg/logger"
	"go.uber.org/zap"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/pkg/waiting"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorGettingBot = "error getting bot"
)

type Command func(ctx context.Context, update tgbotapi.Update) error
type Bot struct {
	bot      *tgbotapi.BotAPI
	mu       sync.Mutex
	upd      tgbotapi.UpdateConfig
	commands map[string]Command
	Waiter   *waiting.Waiter[int64, tgbotapi.Message]
}

func New(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, vanerrors.Wrap(ErrorGettingBot, err)
	}

	u := tgbotapi.NewUpdate(60)

	return &Bot{
		bot:      bot,
		upd:      u,
		commands: map[string]Command{},
		Waiter:   waiting.New[int64, tgbotapi.Message](),
	}, nil
}

func (b *Bot) Init(commands map[string]Command) {
	b.commands = commands
}

func (b *Bot) Run(ctx context.Context) {
	b.mu.Lock()
	defer b.mu.Unlock()
	defer func() {
		if err := recover(); err != nil {
			logger.Error(ctx, "Bot panic", zap.Any("error", err))
		}
	}()

	updates := b.bot.GetUpdatesChan(b.upd)

	for {
		select {
		case <-ctx.Done():
			b.bot.StopReceivingUpdates()
			logger.Debug(ctx, "Bot context done")
			return
		case update := <-updates:

			if update.Message == nil {
				continue
			}

			if c, ok := b.commands[update.Message.Command()]; update.Message.IsCommand() && ok {
				go func() {
					ctx, _ := logger.AddId(ctx, update.FromChat().ID, update.Message.MessageID)
					logger.Debug(ctx, "Bot receives command", zap.String(logger.Command, update.Message.Command()))
					err := c(ctx, update)
					if err != nil {
						logger.Warn(ctx, "Bot error", zap.Any("error", err))
					}
					logger.Debug(ctx, "Bot finished serving message", zap.String(logger.Command, update.Message.Command()))
				}()
				continue
			}

			b.Waiter.Check(update.SentFrom().ID, *update.Message)
		}
	}
}

func (b *Bot) GetUsername() string {
	return b.bot.Self.UserName
}
