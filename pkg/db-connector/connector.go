package db_connector

import (
	"errors"
	"fmt"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

var self IConnector

type Connector struct {
	db          *gorm.DB
	isConnected bool
	conf        *Config
	log         loggerService.ILogger
}

func (c Connector) IsConnected() bool {
	return c.isConnected
}

func (c Connector) GetConnection() *gorm.DB {
	return c.db
}

func (c *Connector) Connect() error {
	if c.IsConnected() {
		return errors.New("db already connected")
	}

	var dsn string
	if os.Getenv("DATABASE_URL") != "" {
		dsn = os.Getenv("DATABASE_URL")
	} else {
		if c.conf == nil {
			return errors.New("db config is empty")
		}
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s",
			c.conf.Host, c.conf.Port, c.conf.User, c.conf.DBName, c.conf.Password,
		)
	}
	println(dsn)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		c.db = &gorm.DB{}
		c.isConnected = false

		return err
	}

	c.db = db
	c.isConnected = true

	return nil
}

func (c *Connector) Disconnect() error {
	if !c.IsConnected() {
		return nil
	}

	err := c.db.Close()
	c.isConnected = false

	return err
}

func MustNewDBConnection(log loggerService.ILogger, conf *Config) IConnector {
	if self != nil {
		log.Fatal("db connector already defined", nil, nil)
	}

	self = &Connector{
		conf: conf,
		log:  log,
	}

	return self
}

func GetDBConnection() (IConnector, error) {
	if self == nil {
		return &Connector{}, errors.New("db connector not defined")
	}

	return self, nil
}
