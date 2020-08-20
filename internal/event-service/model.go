package event_service

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

const (
	monday    base = "1"
	thuesday  base = "2"
	wednesday base = "3"
	thursday  base = "4"
	friday    base = "5"
	saturday  base = "6"
	sunday    base = "7"
)

type base string

type EventsType struct {
	gorm.Model
	EventsTypeName string
	DayOfWeek      []int
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
}

func (d *Events) UnmarshalJSON(data []byte) error {
	type income struct {
		TypeId      int      `json:"type_id"`
		Name        string   `json:"name"`
		Days        []string `json:"days"`
		IsActive    bool     `json:"is_active"`
		Description string   `json:"description"`
		Price       int      `json:"price"`
		StartTime   string   `json:"start_time"`
	}

	var inc income

	if err := json.Unmarshal(data, &inc); err != nil {
		return err
	}

	d.Name = inc.Name
	d.TypeId = inc.TypeId
	d.IsActive = inc.IsActive
	d.Days = fmt.Sprintf("{%s}", strings.Join(inc.Days, ","))
	d.StartTime = inc.StartTime
	d.Price = inc.Price
	d.Description = inc.Description

	return nil
}

func (d Events) MarshalJSON() ([]byte, error) {
	type outcome struct {
		gorm.Model
		Type        EventsType `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
		Name        string     `json:"name"`
		Days        []string   `json:"days"`
		IsActive    bool       `json:"is_active"`
		Description string     `json:"description"`
		Price       int        `json:"price"`
		StartTime   string     `json:"start_time"`
	}

	var days = strings.Replace(d.Days, "{", "", 1)
	days = strings.Replace(days, "}", "", 1)

	var out = outcome{
		Model:       d.Model,
		Name:        d.Name,
		IsActive:    d.IsActive,
		Type:        d.Type,
		Days:        strings.Split(days, ","),
		Description: d.Description,
		Price:       d.Price,
		StartTime:   d.StartTime,
	}

	return json.Marshal(out)
}
