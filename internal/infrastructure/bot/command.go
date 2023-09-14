package bot

import (
	"context"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	CmdStart   = "start"
	CmdHelp    = "help"
	CmdPing    = "ping"
	CmdTest    = "test"
	CmdAddWarn = "warn"
	CmdGPT     = "gpt"
)

func (b *Bot) SendMediaGroup(ctx context.Context, p tg.MediaGroupConfig) ([]tg.Message, error) {
	b.rl.Limit(ctx, p.ChatID)
	return b.bot.SendMediaGroup(p)
}

func (b *Bot) SendMessage(ctx context.Context, p tg.MessageConfig) (tg.Message, error) {
	b.rl.Limit(ctx, p.ChatID)
	return b.bot.Send(p)
}

func (b *Bot) SendVideo(ctx context.Context, p tg.VideoConfig) (tg.Message, error) {
	b.rl.Limit(ctx, p.ChatID)
	return b.bot.Send(p)
}

func (b *Bot) SendAction(ctx context.Context, p tg.ChatActionConfig) (tg.Message, error) {
	b.rl.Limit(ctx, p.ChatID)
	return b.bot.Send(p)
}

func (b *Bot) DeleteMessage(ctx context.Context, p tg.DeleteMessageConfig) (*tg.APIResponse, error) {
	b.rl.Limit(ctx, p.ChatID)
	return b.bot.Request(p)
}

func (b *Bot) EditMessageText(ctx context.Context, p tg.EditMessageTextConfig) (tg.Message, error) {
	b.rl.Limit(ctx, p.ChatID)
	return b.bot.Send(p)
}

func (b *Bot) EditMessageTextTry(ctx context.Context, p tg.EditMessageTextConfig) (tg.Message, error) {
	if ok := b.rl.Available(ctx, p.ChatID); !ok {
		return tg.Message{}, nil
	}
	return b.bot.Send(p)
}
