package password_commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func NewPassword(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		ok, err := service.UserExists(ctx, update.SentFrom().ID)
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Adding password failed with error: %v", err))
		}

		if !ok {
			return b.Send(update.FromChat().ID, update.Message.MessageID, "You don't have an account to add a new password")
		}

		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		params := make([]string, 4)
		order := []string{"company", "username", "password", "your master password"}
		for i, ord := range order {
			err = b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Please enter %s", ord))
			if err != nil {
				return err
			}
			select {
			case <-cancel.Canceled():
				return nil
			case <-ctx.Done():
				return b.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, adding password interrupted")
			case answer := <-wait:
				params[i] = answer.Message.Text
			}
		}
		ok, err = service.CheckUserPassword(ctx, update.SentFrom().ID, params[3])
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Adding password failed with error: %v", err))
		}

		if !ok {
			return b.Send(update.FromChat().ID, update.Message.MessageID, "Wrong master password")
		}

		err = service.NewPassword(ctx, update.SentFrom().ID, params[3], params[2], params[0], params[1])
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Adding password failed with error: %v", err))
		}

		return b.Send(update.FromChat().ID, update.Message.MessageID, "Added password")

	}, "add"
}
