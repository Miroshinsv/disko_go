package handlers

import (
	db "disko/utils"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type Direction struct {
	gorm.Model
	Id          int64
	name        string
	description string
	isSingle    bool
	createdAt   time.Time `json:"-"`
	updatedAt   time.Time `json:"-"`
	deletedAt   time.Time `json:"-"`
}


func GetDirectionById(w http.ResponseWriter, r *http.Request) {
	var events []Events
	db.GetConnect().Find(&events)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func GetAllDirections(w http.ResponseWriter, r *http.Request) {
	db := db.GetConnect()
	var events []Events
	db.Find(&events)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}
