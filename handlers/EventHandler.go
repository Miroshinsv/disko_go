package handlers

import (
	db "disko/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type EventsType struct {
	gorm.Model
	EventsTypeName string
	DayOfWeek      []int
	IsRepeatable   bool
}

type Events struct {
	gorm.Model
	Type     EventsType `gorm:"ForeignKey:TypeId;AssociationForeignKey:id"`
	TypeId   int        `json:"-"`
	Name     string
	Days     []string
	IsActive bool
}

type base string
const (
	monday base = "1"
)

func ActivateEventById(w http.ResponseWriter, r *http.Request) {
	db := db.GetConnect().Debug()
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	var events Events
	db.Preload("Type").Find(&events, i).Update(Events{IsActive: true})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func GetEventById(w http.ResponseWriter, r *http.Request) {
	println("<>")
	db := db.GetConnect().Debug()
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	var events Events
	db.Preload("Type").Find(&events, i)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	db := db.GetConnect().Debug()
	var events []Events
	db.Preload("Type").Find(&events)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	var event Events
	json.NewDecoder(r.Body).Decode(&event)
	err := db.GetConnect().Save(&event)
	if err.Error != nil {
		json.NewEncoder(w).Encode("Can't add event")
	} else {
		json.NewEncoder(w).Encode(event)
	}
}

func UpadteEventById(w http.ResponseWriter, r *http.Request) {
	var nEvent Events
	var event Events
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	json.NewDecoder(r.Body).Decode(&nEvent)
	db.GetConnect().Where(&event, i).Update(nEvent)
	json.NewEncoder(w).Encode(event)
}
