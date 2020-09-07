package models

import (
	"time"

	userService "github.com/Miroshinsv/disko_go/internal/user-service"
)

const (
	PositiveVoice = iota
	NegativeVoice
	NeutralVoice
)

type Vote struct {
	ID        int64              `json:"-"`
	Poll      *Poll              `gorm:"ForeignKey:PollId;AssociationForeignKey:id"`
	PollId    int                `json:"-"`
	User      *userService.Users `gorm:"ForeignKey:UserId;AssociationForeignKey:id"`
	UserId    int                `json:"-"`
	Voice     int
	CreatedAt time.Time `json:"created_at"`
}
