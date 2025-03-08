package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

func Help(b *bot.Bot, _ *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, `Bot commands:

/start - Start the bot
/help - View all commands

Account:

/register - create a new account
/update - update password for account
/remove - remove account

Passwords:

/add - Add new password
/my - View all your passwords
/company - View all your passwords only with concrete company

Hint: after viewing passwords list you can do actions with passwords

Other:
/cancel - Cancel any input
/generate - Generate a random password (it wouldn't be added to your password list)`)
	}, "help"
}

func Start(b *bot.Bot, _ *service.Service) (bot.Command, string) {
	return func(ctx context.Context, update tgbotapi.Update) error {
		return b.SendContext(ctx, update.FromChat().ID, update.Message.MessageID, `This is a password management bot 
 
The creators of the bot do not guarantee the security of your passwords, it is better to use more secure services

/help - View all the commands of the bot`)
	}, "start"
}
