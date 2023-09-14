package bot

import (
	"context"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) CmdTest(ctx context.Context, u tg.Update) (err error) {
	prompt := u.Message.CommandArguments()

	b.SendAction(ctx, tg.NewChatAction(u.Message.Chat.ID, tg.ChatTyping))

	answer, err := b.svc.Ask(ctx, prompt)
	if err != nil {
		return err
	}

	msg := tg.NewMessage(u.Message.Chat.ID, tg.EscapeText(tg.ModeMarkdownV2, answer))
	msg.ReplyToMessageID = u.Message.MessageID
	msg.ParseMode = tg.ModeMarkdownV2

	b.SendMessage(ctx, msg)

	return nil
}

// go func() {
// 	i := 0
// 	t := time.Tick(ProgressUpdateInterval)
// loop:
// 	for {
// 		select {
// 		case <-t:
// 			if i > len(HourglassMsg)-1 {
// 				i = 0
// 			}

// 			// b.bot.Send(tg.NewEditMessageText(u.Message.Chat.ID, msg.MessageID, HourglassMsg[i]))
// 			b.EditMessageTextTry(ctx, tg.NewEditMessageText(u.Message.Chat.ID, msg.MessageID, HourglassMsg[i]))

// 			i++
// 		case <-ctx.Done():
// 			break loop

// 		}
// 	}

// }()
