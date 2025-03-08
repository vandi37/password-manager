package user_commands

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/logger"
)

func UpdateUser(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		if err := b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter your master password"); err != nil {
			return err
		}

		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		var last string

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User cancelled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App cancelled")

			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, changing password interrupted")
		case answer := <-wait:
			last = answer.Text
			_, ok, err := service.CheckUserPassword(ctx, update.SentFrom().ID, answer.Text)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Changing password failed with error: %v", err))
			} else if !ok {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Changing password failed: wrong password")
			}
		}

		if err := b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter your new password"); err != nil {
			return err
		}

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User cancelled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App cancelled")
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, changing password interrupted")
		case answer := <-wait:
			err := service.UpdateUser(ctx, update.SentFrom().ID, answer.Text, last)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Changing password failed with error: %v", err))
			}

			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Password changed. Please store your master password in a safe place.")
		}
	}, "update"
}
