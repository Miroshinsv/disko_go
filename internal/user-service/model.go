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
	VkId       string
	Roles      []*roleService.Roles `gorm:"many2many:users_roles;"`
}

func (u Users) IsAdmin() bool {
	if len(u.Roles) == 0 {
		return false
	}

	for _, v := range u.Roles {
		if v.Admin {
			return true
		}
	}

	return false
}
