package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func Cancel(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		if b.Waiter.Cancel(update.SentFrom().ID) {
			return b.Send(update.FromChat().ID, update.Message.MessageID, "Canceled waiting")
		}
		return b.Send(update.FromChat().ID, update.Message.MessageID, "Nothing to cancel")
	}, "cancel"
}
