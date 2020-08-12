package role_service

import "github.com/jinzhu/gorm"

type Roles struct {
	gorm.Model
	Name string
}
