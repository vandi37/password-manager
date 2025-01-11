package password_commands

import (
	"context"
	"slices"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/postgresql/module"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
	"github.com/vandi37/password-manager/pkg/waiting"
)

func ContinueNoChan(b *bot.Bot, service *service.Service, passwords []module.Password) bot.Command {
	return func(ctx context.Context, update tgbotapi.Update) error {
		wait, cancel := b.Waiter.Add(update.SentFrom().ID)
		defer b.Waiter.Remove(update.SentFrom().ID)

		return Continue(b, service, passwords, wait, cancel)(ctx, update)
	}
}

func Continue(b *bot.Bot, service *service.Service, passwords []module.Password, wait chan tgbotapi.Message, cancel waiting.Cancel) bot.Command {
	return func(ctx context.Context, update tgbotapi.Update) error {

		var n int = -1
		var err error

		select {
		case <-cancel.Canceled():
			return nil
		case <-ctx.Done():
			return b.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, choosing password interrupted")
		case answer := <-wait:
			n, err = strconv.Atoi(answer.Text)
		}

		if err != nil || n < 0 || n > len(passwords) {
			return b.Send(update.FromChat().ID, update.Message.MessageID, "You haven't entered a index that is in range of list of commands")
		}

		password := passwords[n-1]

		actions := []string{"view", "update username", "update password", "remove"}
		commands := []bot.Command{ViewPassword(b, service, password, wait, cancel), UpdatePasswordUsername(b, service, password, wait, cancel), UpdatePassword(b, service, password, wait, cancel), RemovePassword(b, service, password, wait, cancel)}

		err = b.Send(update.FromChat().ID, update.Message.MessageID, "Please choose actions in range of:"+func() string {
			var res string
			for _, a := range actions {
				res += "\n`" + a + "`"
			}
			return res
		}())
		if err != nil {
			return err
		}

		select {
		case <-cancel.Canceled():
			return nil
		case <-ctx.Done():
			return b.Send(update.FromChat().ID, update.Message.MessageID, "I'm sorry, choosing action interrupted")
		case answer := <-wait:
			index := slices.Index(actions, answer.Text)
			if index < 0 {
				return b.Send(update.FromChat().ID, update.Message.MessageID, "You haven't entered a action that was in the list")
			}

			return commands[index](ctx, update)
		}
	}
}
