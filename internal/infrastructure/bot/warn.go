package bot

import (
	"context"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) CmdAddWarn(ctx context.Context, u tg.Update) (err error) {
	msg := tg.NewMessage(u.Message.Chat.ID, "🗿 Добавить предупреждение")
	b.SendMessage(ctx, msg)

	return nil
}

func (b *Bot) CmdMyWarns(ctx context.Context, u tg.Update) (err error) {

	return nil
}

func (b *Bot) CmdGetWarns(ctx context.Context, u tg.Update) (err error) {

	return nil
}
