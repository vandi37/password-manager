package bot

import (
	"context"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/pkg/logger"
	"github.com/vandi37/password-manager/pkg/waiting"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorGettingBot = "error getting bot"
	ContextExit     = "context exit"
)

type Command func(ctx context.Context, update tgbotapi.Update) error

type Bot struct {
	bot      *tgbotapi.BotAPI
	logger   *logger.Logger
	mu       sync.Mutex
	upd      tgbotapi.UpdateConfig
	commands map[string]Command
	Waiter   *waiting.Waiter[int64, tgbotapi.Update]
}

func New(token string, logger *logger.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorGettingBot, err, vanerrors.EmptyHandler)
	}

	u := tgbotapi.NewUpdate(60)

	return &Bot{
		bot:      bot,
		logger:   logger,
		upd:      u,
		commands: map[string]Command{},
		Waiter:   waiting.New[int64, tgbotapi.Update](),
	}, nil
}

func (b *Bot) Init(commands map[string]Command) {
	b.commands = commands
}

func (b *Bot) Run(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	updates := b.bot.GetUpdatesChan(b.upd)

	for {
		select {
		case <-ctx.Done():
			return vanerrors.NewSimple(ContextExit)
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			if c, ok := b.commands[update.Message.Command()]; update.Message.IsCommand() && ok {
				go func() {
					err := c(ctx, update)
					if err != nil {
						b.logger.Errorf("error handling text '%s', user `%d` (%s): %v", update.Message.Text, update.SentFrom().ID, update.SentFrom().FirstName)
					}
				}()
				continue
			}

			b.Waiter.Check(update.SentFrom().ID, update)

		}
	}
}

func (b *Bot) GetUsername() string {
	return b.bot.Self.UserName
}
