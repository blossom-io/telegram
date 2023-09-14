package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Bot         `yaml:"bot"`
		Connections `yaml:"connections"`
		AI          `yaml:"ai"`
	}

	Bot struct {
		Commands struct {
			Enabled struct {
				CmdGPT bool `yaml:"gpt" env:"BLSM_TG_BOT_CMD_GPT"`
			} `yaml:"enabled"`
		} `yaml:"commands"`
		//ApiID       int    `env-required:"true" yaml:"api_id"  env:"BLSM_TG_BOT_API_ID"`
		//ApiHash     string `env-required:"true" yaml:"api_hash" env:"BLSM_TG_BOT_API_HASH"`
		BotToken    string `env-required:"true" yaml:"bot_token" env:"BLSM_TG_BOT_TOKEN"`
		BotUsername string `env-required:"true" yaml:"bot_username" env:"BLSM_TG_BOT_USERNAME"`
		WebHookURL  string `env-required:"true" yaml:"webhook_url" env:"BLSM_TG_WEBHOOK_URL"`
		LogLevel    string `env-required:"true" yaml:"log_level" env:"BLSM_TG_LOG_LEVEL"`
	}

	AI struct {
		OpenAiApiKey       string `env-required:"true" yaml:"openai_api_key" env:"BLSM_TG_AI_API_KEY"`
		MaxTokens          int    `env-required:"true" yaml:"max_tokens" env:"BLSM_TG_AI_MAX_TOKENS"`
		CustomInstructions string `env-required:"false" yaml:"custom_instructions" env:"BLSM_TG_AI_CUSTOM_INSTRUCTIONS"`
	}

	Connections struct {
		Postgres `yaml:"postgres"`
	}

	Postgres struct {
		URL string `env-required:"true" yaml:"url" env:"BLSM_TG_PG_URL"`
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
