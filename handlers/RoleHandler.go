package handlers

import (
	db "disko/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type Roles struct {
	gorm.Model
	Name       string
}

func DisbandRoleById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	db.GetConnect().Where(Roles{}, i).Delete(&Roles{})
	json.NewEncoder(w).Encode("role disband")
}

func AddRole(w http.ResponseWriter, r *http.Request) {
	var role Roles
	json.NewDecoder(r.Body).Decode(&role)
	db.GetConnect().Save(&role)
	json.NewEncoder(w).Encode(role)
}

func UpdateRoleById(w http.ResponseWriter, r *http.Request) {
	var nrole Roles
	var role Roles
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	json.NewDecoder(r.Body).Decode(&nrole)
	db.GetConnect().Where(&role, i).Update(nrole)
	json.NewEncoder(w).Encode(role)
}

func GetRoleById(w http.ResponseWriter, r *http.Request) {
	var role Roles
	i, _ := strconv.Atoi(mux.Vars(r)["id"])
	db.GetConnect().First(&role, i)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(role)
}

func GetAllRoles(w http.ResponseWriter, r *http.Request) {
	db := db.GetConnect()
	var roles []Roles
	db.Find(&roles)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}
