package handlers

import (
	db "disko/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type School struct {
	gorm.Model
	SchoolName  string
	Description string
	Phone       string
	Site        string
	Email       string
	Owner       Users `gorm:"ForeignKey:OwnerId;AssociationForeignKey:id"`
	OwnerId     int
}

func GetAllSchools(w http.ResponseWriter, r *http.Request) {
	var schools []School
	db.GetConnect().Find(&schools)
	json.NewEncoder(w).Encode(schools)
}

func AddSchool(w http.ResponseWriter, r *http.Request) {
	var school School
	json.NewDecoder(r.Body).Decode(&school)
	fmt.Println(school)
	db.GetConnect().Preload("Owner").Save(&school).Find(&school)
	json.NewEncoder(w).Encode(school)
}

func GetSchoolById(w http.ResponseWriter, r *http.Request) {
	var school School
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	db.GetConnect().Preload("Owner").Find(&school, i)
	json.NewEncoder(w).Encode(school)
}

func UpdateSchoolById(w http.ResponseWriter, r *http.Request) {
	var school School
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	json.NewDecoder(r.Body).Decode(&school)
	db.GetConnect().Debug().Update(&school, i)
	json.NewEncoder(w).Encode(school)
}

func DeleteSchoolById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	db.GetConnect().Delete(School{}, i)
	json.NewEncoder(w).Encode("School deleted")
}
