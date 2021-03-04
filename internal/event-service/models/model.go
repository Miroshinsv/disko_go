package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/Miroshinsv/disko_go/internal/poll-service/models"
)

var regReplace = regexp.MustCompile("[{}\"]")

type EventsType struct {
	gorm.Model
	EventsTypeName string
	IsRepeatable   bool
}

type Events struct {
	gorm.Model
	Type        EventsType    `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
	Polls       []models.Poll `gorm:"foreignKey:EventId"`
	TypeId      int           `json:"type_id"`
	Name        string        `json:"name"`
	Days        string        `json:"days"`
	IsActive    bool          `json:"is_active"`
	Description string        `json:"description"`
	Price       string        `json:"price"`
	StartTime   string        `json:"start_time"`
	Logo        string        `json:"logo"`
	Lat         float32       `json:"lat"`
	Lng         float32       `json:"lng"`
	OwnerId     uint          `json:"owner_id"`
}

func (d *Events) UnmarshalJSON(data []byte) error {
	type income struct {
		TypeId      int     `json:"type_id"`
		Name        string  `json:"name"`
		Days        string  `json:"days"`
		IsActive    bool    `json:"is_active"`
		Description string  `json:"description"`
		Price       string  `json:"price"`
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
		Type        EventsType    `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
		Polls       []models.Poll `gorm:"foreignKey:EventId"`
		Name        string        `json:"name"`
		Days        string        `json:"days"`
		IsActive    bool          `json:"is_active"`
		Description string        `json:"description"`
		Price       string        `json:"price"`
		StartTime   string        `json:"start_time"`
		Logo        string        `json:"logo"`
		Lat         float32       `json:"lat"`
		Lng         float32       `json:"lng"`
	}

	for i, j := 0, len(d.Polls)-1; i < j; i, j = i+1, j-1 {
		d.Polls[i], d.Polls[j] = d.Polls[j], d.Polls[i]
	}

	var out = outcome{
		Model:       d.Model,
		Name:        d.Name,
		IsActive:    d.IsActive,
		Type:        d.Type,
		Polls:       d.Polls,
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
