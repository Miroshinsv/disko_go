package user_service

import (
	"github.com/jinzhu/gorm"

	roleService "github.com/Miroshinsv/disko_go/internal/role-service"
)

type Users struct {
	gorm.Model
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Phone      string
	Password   string `json:"-"`
	AvatarUrl  string
	Roles      []*roleService.Roles `gorm:"many2many:users_roles;"`
}
