package handlers

import (
	db "disko/utils"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type EventsType struct {
	gorm.Model
	Id        int `gorm:"primary_key"`
	Name      string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type Events struct {
	gorm.Model
	Id        int64
	Type      EventsType `gorm:"foreignkey:id;association_foreignkey:type_id"`
	TypeId    int
	Name      string
	createdAt time.Time `json:"-"`
	updatedAt time.Time `json:"-"`
	deletedAt time.Time `json:"-"`
}

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	db := db.GetConnect().Debug()
	var events []Events
	db.Set("gorm:auto_preload", true).Find(&events)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}
