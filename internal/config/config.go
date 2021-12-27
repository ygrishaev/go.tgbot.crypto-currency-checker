package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type (
	AppConfig struct {
		TelegramBotName  string `json:"telegram_bot_name" mapstructure:"telegram_bot_name"`
		TelegramAPIToken string `json:"telegram_api_token" mapstructure:"telegram_api_token"`
		CoincapApiUrl    string `json:"coincap_api_url" mapstructure:"coincap_api_url"`
	}
)

func New(configPath string) (*AppConfig, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var appConfig AppConfig
	if err := v.Unmarshal(&appConfig); err != nil {
		return nil, fmt.Errorf("failed to config file: %w", err)
	}

	return &appConfig, nil
}
