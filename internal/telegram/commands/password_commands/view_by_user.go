package password_commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func ViewByUser(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		ok, err := service.UserExists(ctx, update.SentFrom().ID)
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Getting password data failed with error: %v", err))
		}

		if !ok {
			return b.Send(update.FromChat().ID, update.Message.MessageID, "You don't have an account to get password data")
		}

		passwords, err := service.GetPasswordByUserId(ctx, update.SentFrom().ID)
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Getting password data failed with error: %v", err))
		}

		err = b.Send(update.FromChat().ID, update.Message.MessageID, ToString(passwords, ""))
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Getting password data failed with error: %v", err))
		}

		return Continue(b, service, passwords)(ctx, update)
	}, "my"
}
