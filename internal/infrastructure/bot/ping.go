package bot

import (
	"context"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) CmdPing(ctx context.Context, u tg.Update) error {
	if u.Message.Chat.Type != "private" {
		return nil
	}

	b.bot.Send(tg.NewMessage(u.Message.Chat.ID, "pong"))

	return nil
}
