package models

import (
	"time"

	"github.com/jinzhu/gorm"

	event_service "github.com/Miroshinsv/disko_go/internal/event-service"
)

type Poll struct {
	gorm.Model
	Name     string
	Subject  string
	Event    *event_service.Events `gorm:"ForeignKey:EventId;AssociationForeignKey:id"`
	EventId  int                   `json:"-"`
	IsHidden bool
	DueDate  time.Time
}
