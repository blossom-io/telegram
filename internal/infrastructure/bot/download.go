package bot

import (
	"blossom/internal/entity"
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ProgressUpdateInterval = 2 * time.Second
)

var (
	HourglassMsg   = []string{"⏳...", "⌛️.", "⏳.."}
	SupportedSites = []string{
		"tiktok.com",
		"/shorts/",
		"clips.twitch.tv/", "/clip/",
		"instagram.com/reel/",
	}
)

// Downloader checks if the message contains an url and if it is enabled in chat downloads it
func (b *Bot) Downloader(ctx context.Context, u tg.Update) (err error) {
	ctx, stop := context.WithTimeout(ctx, 30*time.Second)
	defer stop()

	urls := collectURLs(u)
	if len(urls) == 0 {

		return nil
	}

	if !isSupported(urls[0], SupportedSites) {

		return nil
	}

	textContainsOnlyURL := len(u.Message.Text) == len(urls[0])

	ok, err := b.svc.IsDownloaderEnabled(ctx, u.Message.Chat.ID)
	if err != nil || !ok {

		return nil
	}

	msg, err := b.SendProgressMsg(ctx, u)
	if err != nil {

		return err
	}
	defer b.DeleteMessage(ctx, tg.NewDeleteMessage(u.Message.Chat.ID, msg.MessageID))

	media, err := b.svc.Fetch(ctx, urls[0])
	if err != nil {

		return err
	}

	m := b.prepareMsg(ctx, u, media)

	b.SendMediaGroup(ctx, m)

	if textContainsOnlyURL {
		b.DeleteMessage(ctx, tg.NewDeleteMessage(u.Message.Chat.ID, u.Message.MessageID))
	}

	return nil
}

// prepareMsg prepares a tg.MediaGroupConfig for sending
func (b *Bot) prepareMsg(ctx context.Context, u tg.Update, media *entity.Media) tg.MediaGroupConfig {
	vid := tg.NewInputMediaVideo(tg.FileReader{
		Name:   "1.mp4",
		Reader: bytes.NewReader(media.Body),
	})
	vid.Thumb = tg.FileReader{
		Name:   "1.jpg",
		Reader: bytes.NewReader(media.Preview),
	}
	vid.Width = media.Width
	vid.Height = media.Height
	vid.Duration = int(media.Duration)
	vid.SupportsStreaming = true
	vid.ParseMode = tg.ModeMarkdownV2
	vid.Caption = fmt.Sprintf("[%s](%s)", media.Extractor, media.WebpageURL)

	m := tg.NewMediaGroup(u.Message.Chat.ID, []any{vid})
	if u.Message.ReplyToMessage != nil {
		m.ReplyToMessageID = u.Message.ReplyToMessage.MessageID
	}

	return m
}

func (b *Bot) SendProgressMsg(ctx context.Context, u tg.Update) (msg tg.Message, err error) {
	reply := tg.NewMessage(u.Message.Chat.ID, "⌛️...")
	reply.ReplyToMessageID = u.Message.MessageID

	msg, err = b.SendMessage(ctx, reply)
	if err != nil {
		return msg, err
	}

	b.SendAction(ctx, tg.NewChatAction(u.Message.Chat.ID, tg.ChatTyping))

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

	return msg, nil
}

func inlineButton(media entity.Media) tg.InlineKeyboardMarkup {
	if media.Title == "" {
		media.Title = media.WebpageURL
	}

	return tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonURL(media.Title, media.WebpageURL),
		),
	)
}

// collectURLs collects urls from the message
func collectURLs(u tg.Update) (urls []string) {
	if u.Message.Entities == nil {

		return nil
	}

	urlPattern := regexp.MustCompile(`https?://[^\s]+`)

	matches := urlPattern.FindAllString(u.Message.Text, -1)

	urls = append(urls, matches...)

	return urls
}

func isSupported(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
