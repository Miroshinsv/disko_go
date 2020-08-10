package db_connector

import "github.com/jinzhu/gorm"

type IConnector interface {
	GetConnection() *gorm.DB
	Connect() error
	Disconnect() error
	IsConnected() bool
}
