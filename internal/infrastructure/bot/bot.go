package bot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"blossom/internal/config"
	"blossom/internal/service"
	"blossom/pkg/logger"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Boter interface {
	Run(ctx context.Context) error
	Close() error
}

type Bot struct {
	log logger.Logger
	cfg *config.Config
	svc service.Servicer
	bot *tg.BotAPI
}

func New(ctx context.Context, cfg *config.Config, log logger.Logger, svc service.Servicer) (Boter, error) {
	bot, err := tg.NewBotAPI(cfg.BotToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Info("Authorized", "username", bot.Self.UserName)

	return &Bot{
		cfg: cfg,
		log: log,
		svc: svc,
		bot: bot,
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	whCfg, err := b.makeWebhookConfig()
	if err != nil {
		return err
	}

	_, err = b.bot.Request(whCfg)
	if err != nil {
		b.log.Error(err.Error())
	}

	updates := b.bot.ListenForWebhook("/" + b.bot.Token)
	go http.ListenAndServe("0.0.0.0:8443", nil)

	for u := range updates {
		if u.Message != nil {
			b.log.Info("msg", "title", u.Message.Chat.Title, "username", u.Message.From.UserName, "text", u.Message.Text)
			b.Downloader(ctx, u)
		}

		if u.Message == nil {
			continue
		}

		if !u.Message.IsCommand() {
			continue
		}

		switch u.Message.Command() {
		case CmdStart:
			err = b.CmdStart(ctx, u)
			if err != nil {
				b.log.Error(err.Error())
				continue
			}
		case CmdHelp:
			err = b.CmdHelp(ctx, u)
			if err != nil {
				b.log.Error(err.Error())
				continue
			}
		case CmdPing:
			err = b.CmdPing(ctx, u)
			if err != nil {
				b.log.Error(err.Error())
				continue
			}
		case CmdTest:
			err = b.CmdTest(ctx, u)
			if err != nil {
				b.log.Error(err.Error())
				continue
			}
		}
	}

	return nil
}

func (b *Bot) Close() error {
	return nil
}

func (b *Bot) makeWebhookConfig() (tg.WebhookConfig, error) {
	URL, err := url.Parse(b.cfg.WebHookURL)
	if err != nil {
		return tg.WebhookConfig{}, err
	}

	whURL := fmt.Sprintf("%s/%s", URL.String(), b.bot.Token)

	b.log.Info("webhook url", "url", whURL)

	return tg.NewWebhook(whURL)
}
