package web_server

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	IsEnabled bool   `mapstructure:"WEB_SERVER_ENABLED"`
	Host      string `mapstructure:"WEB_HOST"`
	Port      int    `mapstructure:"WEB_PORT"`
}

func (c Config) Validate() error {
	if !c.IsEnabled {
		return nil
	}

	if c.Host == "" {
		if os.Getenv("HOST") != "" {
			c.Host = os.Getenv("HOST")
		} else {
			return errors.New("WEB_HOST is not defined")
		}
	}

	c.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	fmt.Println("Port: ")
	fmt.Println(c.Port)
	if c.Port == 0 {
		return errors.New("WEB_PORT is not defined")
	}

	return nil
}
