package password_commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/postgresql/module"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/waiting"
)

func ViewPassword(b *bot.Bot, service *service.Service, password module.Password, wait chan tgbotapi.Message, cancel waiting.Cancel) bot.Command {
	return func(ctx context.Context, update tgbotapi.Update) error {
		err := b.Send(update.FromChat().ID, update.Message.MessageID, "Please enter your master password")
		if err != nil {
			return err
		}

		select {
		case <-cancel.Canceled():
			return nil
		case <-ctx.Done():
			return b.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, viewing password interrupted")
		case answer := <-wait:
			ok, err := service.CheckUserPassword(ctx, update.SentFrom().ID, answer.Text)
			if err != nil {
				return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Viewing password failed with error: %v", err))
			} else if !ok {
				return b.Send(update.FromChat().ID, update.Message.MessageID, "Viewing password failed: wrong password")
			}

			res, err := service.Decrypt([]byte(answer.Text), password.Password, password.Nonce)
			if err != nil {
				return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Viewing password failed with error: %v", err))
			}

			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Password for `%s`:`%s` is `%s`", password.Company, password.Username, res))
		}
	}
}
