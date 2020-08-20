package event_service

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

const (
	monday base = "1"
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
	Type     EventsType `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
	TypeId   int        `json:"type_id"`
	Name     string     `json:"name"`
	Days     string     `json:"days"`
	IsActive bool       `json:"is_active"`
}

func (d *Events) UnmarshalJSON(data []byte) error {
	type income struct {
		TypeId   int      `json:"type_id"`
		Name     string   `json:"name"`
		Days     []string `json:"days"`
		IsActive bool     `json:"is_active"`
	}

	var inc income

	if err := json.Unmarshal(data, &inc); err != nil {
		return err
	}

	d.Name = inc.Name
	d.TypeId = inc.TypeId
	d.IsActive = inc.IsActive
	d.Days = fmt.Sprintf("{%s}", strings.Join(inc.Days, ","))

	return nil
}

func (d Events) MarshalJSON() ([]byte, error) {
	type outcome struct {
		gorm.Model
		Type     EventsType `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
		Name     string     `json:"name"`
		Days     []string   `json:"days"`
		IsActive bool       `json:"is_active"`
	}

	var days = strings.Replace(d.Days, "{", "", 1)
	days = strings.Replace(days, "}", "", 1)

	var out = outcome{
		Model:    d.Model,
		Name:     d.Name,
		IsActive: d.IsActive,
		Type:     d.Type,
		Days:     strings.Split(days, ","),
	}

	return json.Marshal(out)
}
