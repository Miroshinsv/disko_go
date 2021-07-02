package role_service

import "github.com/jinzhu/gorm"

type Roles struct {
	gorm.Model
	Name        string
	Admin       bool `json:"admin"`
	SchoolAdmin bool `json:"school_admin"`
	Dj          bool `json:"dj"`
}
