package school_service

import (
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	"github.com/jinzhu/gorm"
)

type School struct {
	gorm.Model
	SchoolName  string
	Description string
	Phone       string
	Site        string
	Email       string
	Owner       userService.Users `gorm:"ForeignKey:OwnerId;AssociationForeignKey:id"`
	OwnerId     int
}
