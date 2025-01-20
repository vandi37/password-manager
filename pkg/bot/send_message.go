package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorSending = "error sending"
	TextEmpty    = "text empty"
)

func (b *Bot) Send(chat int64, repl int, text string) error {
	if text == "" {
		return vanerrors.Simple(TextEmpty)
	}

	messageParts := strings.Split(text, "\n")

	var sendText string
	last := len(messageParts) - 1

	for i := 0; i < len(messageParts); i++ {
		part := messageParts[i]
		length := len(part) + len(sendText)
		if length >= 4000 || last == i {

			if last == i {
				sendText += "\n" + part
			}

			msg := tgbotapi.NewMessage(chat, sendText)
			msg.ParseMode = "Markdown"
			msg.ReplyToMessageID = repl
			msg.DisableWebPagePreview = true

			_, err := b.bot.Send(msg)
			if err != nil {
				return vanerrors.Wrap(ErrorSending, err)
			}

			sendText = part
		} else {
			sendText += "\n" + part
		}
	}
	return nil
}
