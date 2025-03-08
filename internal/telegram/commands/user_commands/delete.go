package user_commands

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/logger"
)

func DeleteUser(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		ok, err := service.UserExists(ctx, update.SentFrom().ID)
		if err != nil {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Removing account failed with error: %v", err))
		}

		if !ok {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "You don't have an account to remove")
		}

		if err := b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Are you sure, that you want to remove your account? You can't recover your data after deleting, all your passwords would be removed forever\n\n Please write `yes` if you really want"); err != nil {
			return err
		}

		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		answers := []string{"yes", "i'm sure", "definitely"}
		for i, a := range answers {
			select {
			case <-cancel.Canceled():
				logger.Debug(ctx, "User cancelled")
				return nil
			case <-ctx.Done():
				logger.Debug(ctx, "App cancelled")
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "I'm sorry, removing account interrupted")
			case answer := <-wait:
				if answer.Text != a {
					return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("You've entered `%s`, not `%s`, removing account failed", answer.Text, a))
				}
			}
			if i+1 < len(answers) {
				return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Now please write `%s`", answers[i+1]))
			}
		}
		err = service.RemoveUser(ctx, update.SentFrom().ID)
		if err != nil {
			return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Removing account failed with error: %v", err))
		}

		return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, "Removed your account")
	}, "remove"
}
