package user_service

import "github.com/jinzhu/gorm"

type Users struct {
	gorm.Model
	FirstName  string
	SureName   string
	MiddleName string
	Email      string
	Phone      string
	Password   string `json:"-"`
}
