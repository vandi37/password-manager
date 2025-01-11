package user_commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func DeleteUser(b *bot.Bot, service *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		ok, err := service.UserExists(ctx, update.SentFrom().ID)
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Removing account failed with error: %v", err))
		}

		if !ok {
			return b.Send(update.FromChat().ID, update.Message.MessageID, "You don't have an account to remove")
		}

		err = b.Send(update.FromChat().ID, update.Message.MessageID, "Are you sure, that you want to remove your account? You can't recover your data after deleting, all your passwords would be removed forever\n\n Please write `yes` if you really want")
		if err != nil {
			return err
		}

		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		answers := []string{"yes", "i'm sure", "definitely"}
		for i, a := range answers {
			select {
			case <-cancel.Canceled():
				return nil
			case <-ctx.Done():
				return b.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, removing account interrupted")
			case answer := <-wait:
				if answer.Text != a {
					return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("You've entered `%s`, not `%s`, removing account failed", answer.Text, a))
				}
			}
			if i+1 < len(answers) {
				err := b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Now please write `%s`", answers[i+1]))
				if err != nil {
					return err
				}
			}
		}
		err = service.RemoveUser(ctx, update.SentFrom().ID)
		if err != nil {
			return b.Send(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Removing account failed with error: %v", err))
		}

		return b.Send(update.FromChat().ID, update.Message.MessageID, "Removed your account")
	}, "remove"
}
