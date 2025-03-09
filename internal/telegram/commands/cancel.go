package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func Cancel(b *bot.Bot, _ *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		if b.Waiter.Remove(update.SentFrom().ID) {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Canceled waiting")
		}
		return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Nothing to cancel")
	}, "cancel"
}
