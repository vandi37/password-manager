package commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func NewUser(bot *bot.Bot, service *service.Service) bot.Command {
	return func(ctx context.Context, update tgbotapi.Update) error {
		err := bot.Send(update.FromChat().ID, update.Message.MessageID, "Please enter your master password", nil)
		if err != nil {
			return err
		}

		wait := bot.Waiter.Add(update.SentFrom().ID)
		defer bot.Waiter.Remove(update.SentFrom().ID)

		select {
		case <-ctx.Done():

			return bot.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, registration interrupted", nil)

		case answer := <-wait:
			err := service.NewUser(ctx, update.SentFrom().ID, answer.Message.Text)
			if err != nil {
				return bot.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Registration failed with error: %v", err), nil)
			} else {
				return bot.Send(update.FromChat().ID, update.Message.MessageID, "Registration finished. Please store your master password in a safe place.", nil)
			}
		}
	}
}
