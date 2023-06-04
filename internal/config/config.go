package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Bot         `yaml:"bot"`
		Connections `yaml:"connections"`
	}

	Bot struct {
		//ApiID       int    `env-required:"true" yaml:"api_id"  env:"BLSM_TG_BOT_API_ID"`
		//ApiHash     string `env-required:"true" yaml:"api_hash" env:"BLSM_TG_BOT_API_HASH"`
		BotToken    string `env-required:"true" yaml:"bot_token" env:"BLSM_TG_BOT_TOKEN"`
		BotUsername string `env-required:"true" yaml:"bot_username" env:"BLSM_TG_BOT_USERNAME"`
		WebHookURL  string `env-required:"true" yaml:"webhook_url" env:"BLSM_TG_WEBHOOK_URL"`
		LogLevel    string `env-required:"true" yaml:"log_level" env:"BLSM_TG_LOG_LEVEL"`
	}

	Connections struct {
		Postgres `yaml:"postgres"`
	}

	Postgres struct {
		URL string `env-required:"true" yaml:"url" env:"PG_URL"`
	}
)

// New returns app config.
func New() (*Config, error) {
	c := &Config{}

	err := cleanenv.ReadEnv(c)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return c, nil
}
