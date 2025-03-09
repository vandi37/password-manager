package user_commands

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/logger"
	"go.uber.org/zap"
)

func NewUser(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		ok, err := service.UserExists(ctx, update.SentFrom().ID)
		if err != nil {
			logger.Warn(ctx, "UserExists error", zap.Error(err))
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Registration failed with error")
		}
		if ok {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "User already exists")
		}
		err = b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter your master password")
		if err != nil {
			return err
		}

		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User cancelled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App cancelled")

			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, registration interrupted")
		case answer := <-wait:
			err := service.NewUser(ctx, update.SentFrom().ID, answer.Text)
			if err != nil {
				logger.Warn(ctx, "Failed registration", zap.Error(err))
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Registration failed with error")
			}
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Registration finished. Please store your master password in a safe place.")
		}
	}, "register"
}
