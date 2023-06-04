package bot

import (
	"context"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) CmdHelp(ctx context.Context, u tg.Update) error {
	return nil
}
