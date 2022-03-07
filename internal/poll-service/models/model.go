package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Poll struct {
	gorm.Model
	Name    string
	Subject string
	// Event    *event_service.Events `gorm:"ForeignKey:EventId;AssociationForeignKey:id"`
	EventId   int `json:"-"`
	IsHidden  bool
	VoteCount int
	DueDate   time.Time
}
