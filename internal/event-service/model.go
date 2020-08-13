package event_service

import "github.com/jinzhu/gorm"

const (
	monday base = "1"
)

type base string

type EventsType struct {
	gorm.Model
	EventsTypeName string
	DayOfWeek      int
	IsRepeatable   bool
}

type Events struct {
	gorm.Model
	Type     EventsType `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
	TypeId   int        `json:"-"`
	Name     string
	Days     []string
	IsActive bool
}
