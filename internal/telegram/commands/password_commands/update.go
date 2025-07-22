package password_commands

import (
	"context"
	"fmt"
	"github.com/vandi37/password-manager/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/postgresql/module"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/waiting"
)

func UpdatePasswordUsername(b *bot.Bot, service *service.Service, password module.Password, wait chan tgbotapi.Message, cancel waiting.Cancel) bot.Command {
	return func(ctx context.Context, update tgbotapi.Update) error {
		err := b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter your master password")
		if err != nil {
			return err
		}

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User canceled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App canceled")
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, updating password username interrupted")
		case answer := <-wait:
			_, ok, err := service.CheckUserPassword(ctx, update.SentFrom().ID, answer.Text)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Updating password username failed with error: %v", err))
			} else if !ok {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Updating password username failed: wrong password")
			}
		}

		err = b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter new username")
		if err != nil {
			return err
		}

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User canceled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App canceled")
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, updating password username interrupted")
		case answer := <-wait:
			err = service.UpdatePasswordUsername(ctx, password.Id, answer.Text)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Updating password username failed with error: %v", err))
			}
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Username updated")
		}

	}
}

func UpdatePassword(b *bot.Bot, service *service.Service, password module.Password, wait chan tgbotapi.Message, cancel waiting.Cancel) bot.Command {
	return func(ctx context.Context, update tgbotapi.Update) error {
		err := b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter your master password")
		if err != nil {
			return err
		}

		var key string

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User canceled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App canceled")
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, updating password interrupted")
		case answer := <-wait:
			var ok bool
			key, ok, err = service.CheckUserPassword(ctx, update.SentFrom().ID, answer.Text)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Updating password failed with error: %v", err))
			} else if !ok {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Updating password failed: wrong password")
			}
		}

		err = b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter new password")
		if err != nil {
			return err
		}

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User canceled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App canceled")
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, updating password interrupted")
		case answer := <-wait:
			err = service.UpdatePassword(ctx, password.Id, answer.Text, key)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Updating password failed with error: %v", err))
			}
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Password updated")
		}

	}
}
