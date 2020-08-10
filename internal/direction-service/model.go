package direction_service

import "github.com/jinzhu/gorm"

type Direction struct {
	gorm.Model
	Name     string
	IsSingle bool
}
