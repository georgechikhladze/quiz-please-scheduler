package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

type Config struct {
	Schedule string         `yaml:"schedule"`
	Telegram TelegramConfig `yaml:"telegram"`
	HTTPPort int            `yaml:"http_port"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
