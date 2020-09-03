package event_service

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"regexp"
	"strings"
)

var regReplace = regexp.MustCompile("[{}\"]")

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
	Lat         float32    `json:"lat"`
	Lng         float32    `json:"lng"`
}

func (d *Events) UnmarshalJSON(data []byte) error {
	type income struct {
		TypeId      int     `json:"type_id"`
		Name        string  `json:"name"`
		Days        string  `json:"days"`
		IsActive    bool    `json:"is_active"`
		Description string  `json:"description"`
		Price       int     `json:"price"`
		StartTime   string  `json:"start_time"`
		Logo        string  `json:"logo"`
		Lat         float32 `json:"lat"`
		Lng         float32 `json:"lng"`
	}

	var inc income

	if err := json.Unmarshal(data, &inc); err != nil {
		return err
	}

	d.Name = inc.Name
	d.TypeId = inc.TypeId
	d.IsActive = inc.IsActive

	daysSlice := strings.Split(inc.Days, ",")
	d.Days = fmt.Sprintf("{\"%s\"}", strings.Join(daysSlice, "\",\""))

	d.StartTime = inc.StartTime
	d.Price = inc.Price
	d.Description = inc.Description
	d.Logo = inc.Logo
	d.Lat = inc.Lat
	d.Lng = inc.Lng
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
		Lat         float32    `json:"lat"`
		Lng         float32    `json:"lng"`
	}

	var out = outcome{
		Model:       d.Model,
		Name:        d.Name,
		IsActive:    d.IsActive,
		Type:        d.Type,
		Days:        regReplace.ReplaceAllString(d.Days, ""),
		Description: d.Description,
		Price:       d.Price,
		StartTime:   d.StartTime,
		Logo:        d.Logo,
		Lat:         d.Lat,
		Lng:         d.Lng,
	}

	return json.Marshal(out)
}
