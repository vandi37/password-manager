package commands

import (
	"github.com/vandi37/password-manager/internal/service"
	"github.com/vandi37/password-manager/pkg/bot"
)

type CommandBuilder func(*bot.Bot, *service.Service) (bot.Command, string)

func BuildCommands(b *bot.Bot, service *service.Service, cmds ...CommandBuilder) map[string]bot.Command {
	res := map[string]bot.Command{}

	for _, c := range cmds {
		command, key := c(b, service)
		res[key] = command
	}
	return res
}
