package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vandi37/password-manager/pkg/logger"
	"github.com/vandi37/password-manager/pkg/utils"
	"github.com/vandi37/vanerrors"
	"go.uber.org/zap"
	"strings"
)

const (
	ErrorSending  = "error sending"
	TextEmpty     = "text empty"
	TextIsTooLong = "text is too long"
	TooLongText   = 4000
	Markdown      = "Markdown"
)

func (b *Bot) sendShort(chat int64, repl int, text string) error {
	if strings.TrimSpace(text) == "" {
		return vanerrors.Simple(TextEmpty)
	}

	if len(strings.TrimSpace(text)) > TooLongText {
		return vanerrors.Simple(TextIsTooLong)
	}

	msg := tgbotapi.NewMessage(chat, strings.TrimSpace(text))
	msg.ParseMode = Markdown
	if repl > 0 {
		msg.ReplyToMessageID = repl
	}
	msg.DisableWebPagePreview = true
	_, err := b.bot.Send(msg)
	if err != nil {
		return vanerrors.Wrap(ErrorSending, err)
	}
	return nil
}

func (b *Bot) send(chat int64, repl int, text string) error {
	if strings.TrimSpace(text) == "" {
		return vanerrors.Simple(TextEmpty)
	}

	if len(strings.TrimSpace(text)) <= TooLongText {
		return b.sendShort(chat, repl, text)
	}

	parts := strings.Split(text, "\n")
	sendText := ""

	for _, part := range parts {
		if len(part) > TooLongText {
			messages, left := utils.SplitString(part, TooLongText)
			if strings.TrimSpace(sendText) != "" {
				err := b.sendShort(chat, repl, sendText)
				if err != nil {
					return err
				}
			}
			for _, message := range messages {
				if err := b.sendShort(chat, repl, message); err != nil {
					return err
				}
			}
			sendText = left
		} else if len(sendText+"\n"+part) <= TooLongText {
			sendText += "\n" + part
		} else if len(sendText+"\n"+part) > TooLongText {
			if strings.TrimSpace(sendText) != "" {
				err := b.sendShort(chat, repl, sendText)
				if err != nil {
					return err
				}
			}
			sendText = part
		}

	}

	if strings.TrimSpace(sendText) != "" {
		return b.sendShort(chat, repl, sendText)
	}
	return nil
}

func (b *Bot) SendContext(ctx context.Context, chat int64, repl int, text string) error {
	if err := b.send(chat, repl, text); err != nil {
		logger.Debug(ctx, "Failed to send message", zap.Error(err))
		return err
	}
	return nil
}
