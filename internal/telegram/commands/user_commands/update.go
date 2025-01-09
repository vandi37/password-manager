package user_commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func UpdateUser(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		err := b.Send(update.FromChat().ID, update.Message.MessageID, "Please enter your master password")
		if err != nil {
			return err
		}

		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		select {
		case <-cancel.Canceled():
			return nil
		case <-ctx.Done():
			return b.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, changing password interrupted")
		case answer := <-wait:
			ok, err := service.CheckUserPassword(ctx, update.SentFrom().ID, answer.Message.Text)
			if err != nil {
				return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Changing password failed with error: %v", err))
			} else if !ok {
				return b.Send(update.FromChat().ID, update.Message.MessageID, "Changing password failed: wrong password")
			}

		}

		err = b.Send(update.FromChat().ID, update.Message.MessageID, "Please enter your new password")
		if err != nil {
			return err
		}
		select {
		case <-cancel.Canceled():
			return nil
		case <-ctx.Done():
			return b.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, changing password interrupted")
		case answer := <-wait:
			err = service.UpdateUser(ctx, update.SentFrom().ID, answer.Message.Text)
			if err != nil {
				return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Changing password failed with error: %v", err))
			}

			return b.Send(update.FromChat().ID, update.Message.MessageID, "Password changed. Please store your master password in a safe place.")
		}
	}, "change_password"
}
