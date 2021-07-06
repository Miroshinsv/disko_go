package models

import "github.com/jinzhu/gorm"

type City struct {
	gorm.Model
	CityName string `json:"city_name"`
	Country  string `json:"country"`
}

func (City) TableName() string {
	return "city"
}
