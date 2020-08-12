package config_service

import (
	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var self IConfig = nil

type Config struct {
	data     map[string]interface{}
	isLoaded bool
}

func (c *Config) Load(path string) error {
	if c.isLoaded {
		return errors.New("config already loaded")
	}

	_, err := toml.DecodeFile(path, &c.data)
	if err != nil {
		return err
	}

	c.isLoaded = true

	return nil
}

func (c Config) Convert(subject ISubjectInterface) error {
	if !c.isLoaded {
		return errors.New("config not loaded")
	}

	err := mapstructure.Decode(c.data, subject)
	if err != nil {
		return err
	}

	err = subject.Validate()

	return err
}

func GetConfigService() IConfig {
	if self == nil {
		self = &Config{}
	}

	return self
}
