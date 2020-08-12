package db_connector

import "errors"

type Config struct {
	IsEnabled bool   `mapstructure:"DB_ENABLED"`
	Host      string `mapstructure:"DB_HOST"`
	User      string `mapstructure:"DB_USER"`
	DBName    string `mapstructure:"DB_NAME"`
	Password  string `mapstructure:"DB_PASSWORD"`
	Port      int    `mapstructure:"DB_PORT"`
}

func (c Config) Validate() error {
	if !c.IsEnabled {
		return nil
	}

	if c.Host == "" {
		return errors.New("invalid or empty parameter DB_HOST")
	}

	if c.User == "" {
		return errors.New("invalid or empty parameter DB_USER")
	}

	if c.DBName == "" {
		return errors.New("invalid or empty parameter DB_NAME")
	}

	if c.Password == "" {
		return errors.New("invalid or empty parameter DB_PASSWORD")
	}

	if c.Port == 0 {
		return errors.New("invalid or empty parameter DB_PORT")
	}

	return nil
}
