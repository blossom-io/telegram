package bot

import (
	"blossom/internal/entity"
	"bytes"
	"context"
	"regexp"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ProgressUpdateInterval = 2 * time.Second
)

var (
	ProgressMsg    = []string{"⌛️.", "⏳..", "⌛️...", "⏳..."}
	SupportedSites = []string{
		"tiktok.com",
		"/shorts/",
		"clips.twitch.tv/", "/clip/",
		"instagram.com/reel/",
	}
)

// Downloader checks if the message contains an url and if it is enabled in chat downloads it
func (b *Bot) Downloader(ctx context.Context, u tg.Update) (err error) {
	ctx, stop := context.WithTimeout(ctx, 20*time.Second)
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
	defer b.bot.Request(tg.NewDeleteMessage(u.Message.Chat.ID, msg.MessageID))

	media, body, thumb, err := b.svc.Download(ctx, urls[0])
	if err != nil {

		return err
	}

	video := tg.NewVideo(u.Message.Chat.ID, tg.FileReader{
		Name:   "1.mp4",
		Reader: bytes.NewReader(body),
	})
	video.ReplyMarkup = inlineButton(media)
	video.Thumb = tg.FileReader{
		Name:   "1.jpg",
		Reader: bytes.NewReader(thumb),
	}
	video.SupportsStreaming = true
	video.Duration = int(media.Duration)
	if !textContainsOnlyURL {
		video.ReplyToMessageID = u.Message.MessageID
	}

	b.bot.Send(tg.NewChatAction(u.Message.Chat.ID, tg.ChatUploadVideo))

	_, err = b.bot.Send(video)
	if err != nil {

		return err
	}

	if textContainsOnlyURL {
		b.bot.Request(tg.NewDeleteMessage(u.Message.Chat.ID, u.Message.MessageID))
	}

	return nil
}

func (b *Bot) SendProgressMsg(ctx context.Context, u tg.Update) (msg tg.Message, err error) {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	// defer stop()

	reply := tg.NewMessage(u.Message.Chat.ID, "⏳...")
	reply.ReplyToMessageID = u.Message.MessageID

	msg, err = b.bot.Send(reply)
	if err != nil {
		return msg, err
	}

	b.bot.Send(tg.NewChatAction(u.Message.Chat.ID, tg.ChatTyping))

	go func() {
		i := 0
		t := time.Tick(ProgressUpdateInterval)
	loop:
		for {
			select {
			case <-t:
				if i > len(ProgressMsg)-1 {
					i = 0
				}

				b.bot.Send(tg.NewEditMessageText(u.Message.Chat.ID, msg.MessageID, ProgressMsg[i]))

				i++
			case <-ctx.Done():
				break loop

			}
		}

	}()

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
