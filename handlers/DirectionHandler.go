package handlers

import (
	db "disko/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type Direction struct {
	gorm.Model
	Name     string
	IsSingle bool
}

func DisbandDirectionById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	db.GetConnect().Where(Direction{}, i).Delete(&Direction{})
	json.NewEncoder(w).Encode("direction deleted")
}

func AddDirection(w http.ResponseWriter, r *http.Request) {
	var d Direction
	json.NewDecoder(r.Body).Decode(&d)
	err := db.GetConnect().Save(&d).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode(d)
}

func UpdateDirectionById(w http.ResponseWriter, r *http.Request) {
	var newDirection Direction
	var oldDirection Direction
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	json.NewDecoder(r.Body).Decode(&newDirection)
	db.GetConnect().Where(&oldDirection, i).Update(newDirection)
	json.NewEncoder(w).Encode(oldDirection)
}

func GetDirectionById(w http.ResponseWriter, r *http.Request) {
	var direction Direction
	i, _ := strconv.Atoi(mux.Vars(r)["id"])
	db.GetConnect().First(&direction, i)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(direction)
}

func GetAllDirections(w http.ResponseWriter, r *http.Request) {
	var directions []Direction
	db.GetConnect().Find(&directions)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(directions)
}
