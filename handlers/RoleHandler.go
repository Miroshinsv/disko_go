package handlers

import (
	db "disko/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type Roles struct {
	Id         int
	Name       string
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
	DeletedAt  time.Time `json:"-"`
}

func DisbandRoleById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid id")
		return
	}
	db.GetConnect().Where(Roles{Id: i}).Delete(&Roles{})
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
	nrole.ModifiedAt = time.Now()
	db.GetConnect().Where(&role, Roles{Id: i}).Update(nrole)
	json.NewEncoder(w).Encode(role)
}

func GetRoleById(w http.ResponseWriter, r *http.Request) {
	var role Roles
	i, _ := strconv.Atoi(mux.Vars(r)["id"])
	db.GetConnect().First(&role, Roles{Id: i})
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
