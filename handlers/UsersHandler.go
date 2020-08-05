package handlers

import (
	db "disko/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type Users struct {
	gorm.Model
	FirstName  string
	SureName   string
	MiddleName string
	Email      string
	Phone      string
	Password   string `json:"-"`
}

func DisbandUserById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	db.GetConnect().Where(Users{}, i).Delete(&Users{})
	json.NewEncoder(w).Encode("user disband")
}

func AddUser(w http.ResponseWriter, r *http.Request) {
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

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	var newUser Users
	var oldUser Users
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	json.NewDecoder(r.Body).Decode(&newUser)
	db.GetConnect().Where(&oldUser, i).Update(newUser)
	json.NewEncoder(w).Encode(oldUser)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	var user Users
	i, _ := strconv.Atoi(mux.Vars(r)["id"])
	db.GetConnect().First(&user, i)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []Users
	db.GetConnect().Find(&users)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
