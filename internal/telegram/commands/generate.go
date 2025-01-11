package commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/generate"
)

func GeneratePassword(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Your password: %s", generate.GeneratePassword(20, true, true, true, false)))
	}, "generate"
}
