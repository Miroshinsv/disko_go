package auth_service

import (
	"errors"
)

var (
	errorEmptyClientID = errors.New("VK_CLIENT_ID is empty")
	errorEmptySecretID = errors.New("VK_CLIENT_SECRET is empty")
)

type Config struct {
	VKClientID     string `mapstructure:"VK_CLIENT_ID"`
	VKClientSecret string `mapstructure:"VK_CLIENT_SECRET"`
}

func (c Config) Validate() error {
	if c.VKClientID == "" {
		return errorEmptyClientID
	}

	if c.VKClientSecret == "" {
		return errorEmptySecretID
	}

	return nil
}
