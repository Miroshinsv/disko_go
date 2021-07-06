package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"regexp"

	"github.com/Miroshinsv/disko_go/internal/poll-service/models"
)

var regReplace = regexp.MustCompile("[{}\"]")

type EventsType struct {
	gorm.Model
	EventsTypeName string
	IsRepeatable   bool
}

type City struct {
	gorm.Model
	CityName string `json:"city_name"`
	Country  string `json:"country"`
}

func (City) TableName() string {
	return "city"
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
	City        City
	CityID      int
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
		City        City    `json:"city_id"`
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
	d.Lat = inc.Lat
	d.Lng = inc.Lng
	d.City = inc.City
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
		City        City          `gorm:"ForeignKey:CityId"`
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
		City:        d.City,
	}

	return json.Marshal(out)
}
