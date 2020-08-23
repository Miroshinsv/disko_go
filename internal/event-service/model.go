package event_service

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type EventsType struct {
	gorm.Model
	EventsTypeName string
	IsRepeatable   bool
}

type Events struct {
	gorm.Model
	Type        EventsType `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
	TypeId      int        `json:"type_id"`
	Name        string     `json:"name"`
	Days        string     `json:"days"`
	IsActive    bool       `json:"is_active"`
	Description string     `json:"description"`
	Price       int        `json:"price"`
	StartTime   string     `json:"start_time"`
	Logo        string     `json:"logo"`
}

func (d *Events) UnmarshalJSON(data []byte) error {
	type income struct {
		TypeId      int    `json:"type_id"`
		Name        string `json:"name"`
		Days        string `json:"days"`
		IsActive    bool   `json:"is_active"`
		Description string `json:"description"`
		Price       int    `json:"price"`
		StartTime   string `json:"start_time"`
		Logo        string `json:"logo"`
	}

	var inc income

	if err := json.Unmarshal(data, &inc); err != nil {
		return err
	}

	d.Name = inc.Name
	d.TypeId = inc.TypeId
	d.IsActive = inc.IsActive
	d.Days = inc.Days
	d.StartTime = inc.StartTime
	d.Price = inc.Price
	d.Description = inc.Description
	d.Logo = inc.Logo
	return nil
}

func (d Events) MarshalJSON() ([]byte, error) {
	type outcome struct {
		gorm.Model
		Type        EventsType `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
		Name        string     `json:"name"`
		Days        string     `json:"days"`
		IsActive    bool       `json:"is_active"`
		Description string     `json:"description"`
		Price       int        `json:"price"`
		StartTime   string     `json:"start_time"`
		Logo        string     `json:"logo"`
	}

	var out = outcome{
		Model:       d.Model,
		Name:        d.Name,
		IsActive:    d.IsActive,
		Type:        d.Type,
		Days:        d.Days,
		Description: d.Description,
		Price:       d.Price,
		StartTime:   d.StartTime,
		Logo:        d.Logo,
	}

	return json.Marshal(out)
}
