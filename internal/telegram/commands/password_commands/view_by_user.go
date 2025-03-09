package password_commands

import (
	"context"
	"github.com/vandi37/password-manager/pkg/logger"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func ViewByUser(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		ok, err := service.UserExists(ctx, update.SentFrom().ID)
		if err != nil {
			logger.Warn(ctx, "UserExists error", zap.Error(err))
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Getting password data failed with error")
		}

		if !ok {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "You don't have an account to get password data")
		}

		passwords, err := service.GetPasswordsByUserId(ctx, update.SentFrom().ID)
		if err != nil {
			logger.Warn(ctx, "GetPasswordsByUserId error", zap.Error(err))
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Getting password data failed with error")
		}
		mes, ok := ToString(passwords, "")
		err = b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, mes)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		return ContinueNoChan(b, service, passwords)(ctx, update)
	}, "my"
}
