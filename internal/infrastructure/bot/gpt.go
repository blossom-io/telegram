package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) CmdGPT(u tg.Update) (err error) {
	if !b.cfg.Bot.Commands.Enabled.CmdGPT {
		return nil
	}

	if len(u.Message.CommandArguments()) == 0 {
		return nil
	}

	ctx, stop := context.WithTimeout(context.Background(), 20*time.Second)
	defer stop()

	firstMsg := tg.NewMessage(u.Message.Chat.ID, "✍️ ...")
	firstMsg.ReplyToMessageID = u.Message.MessageID

	msg, err := b.SendMessage(ctx, firstMsg)
	if err != nil {
		return err
	}

	b.SendAction(ctx, tg.NewChatAction(u.Message.Chat.ID, tg.ChatTyping))

	prompt := fmt.Sprintf("%s. %s", u.Message.CommandArguments(), b.cfg.AI.CustomInstructions)

	delta, err := b.svc.AskStream(ctx, prompt)
	if err != nil {
		return err
	}
	defer close(delta)

	var text strings.Builder

	for chunk := range delta {
		if chunk.Err != nil {

			break
		}

		text.WriteString(chunk.Content)

		fmt.Printf("%s", chunk.Content)

		txt := fmt.Sprintf("%s %s", text.String(), "✍️...")

		b.EditMessageTextTry(ctx, tg.NewEditMessageText(u.Message.Chat.ID, msg.MessageID, txt))

	}

	b.EditMessageText(context.Background(), tg.NewEditMessageText(u.Message.Chat.ID, msg.MessageID, text.String()))

	<-ctx.Done()

	return nil
}
