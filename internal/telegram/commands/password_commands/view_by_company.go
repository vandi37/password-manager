package password_commands

import (
	"context"
	"fmt"
	"github.com/vandi37/password-manager/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func ViewByCompany(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		ok, err := service.UserExists(ctx, update.SentFrom().ID)
		if err != nil {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Getting password data failed with error: %v", err))
		}

		if !ok {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "You don't have an account to get password data")
		}

		err = b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Please enter company name")
		if err != nil {
			return err
		}

		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		select {
		case <-cancel.Canceled():
			logger.Debug(ctx, "User canceled")
			return nil
		case <-ctx.Done():
			logger.Debug(ctx, "App canceled")
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, choosing password interrupted")
		case answer := <-wait:
			passwords, err := service.GetPasswordsByCompany(ctx, update.SentFrom().ID, answer.Text)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Getting password data failed with error: %v", err))
			}

			mes, ok := ToString(passwords, fmt.Sprintf(" with company %s", answer.Text))
			err = b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, mes)
			if err != nil {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Getting password data failed with error: %v", err))
			}

			if !ok {
				return nil
			}

			return Continue(b, service, passwords, wait, cancel)(ctx, update)
		}

	}, "company"
}
