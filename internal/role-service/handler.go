package role_service

import (
	"encoding/json"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	log  loggerService.ILogger
	conn dbConnector.IConnector
}

func (h Handler) DisbandRoleById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	h.conn.GetConnection().Where(Roles{}, i).Delete(&Roles{})

	_ = json.NewEncoder(w).Encode("role disband")
}

func (h Handler) AddRole(w http.ResponseWriter, r *http.Request) {
	var role Roles
	//@todo: cover error
	_ = json.NewDecoder(r.Body).Decode(&role)

	h.conn.GetConnection().Save(&role)

	_ = json.NewEncoder(w).Encode(role)
}

func (h Handler) UpdateRoleById(w http.ResponseWriter, r *http.Request) {
	var (
		nRole Roles
	)
	var role Roles
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	//@todo: cover error
	_ = json.NewDecoder(r.Body).Decode(&nRole)

	h.conn.GetConnection().Where(&role, i).Update(nRole)

	_ = json.NewEncoder(w).Encode(role)
}

func (h Handler) GetRoleById(w http.ResponseWriter, r *http.Request) {
	var role Roles
	//@todo: cover error
	i, _ := strconv.Atoi(mux.Vars(r)["id"])

	h.conn.GetConnection().First(&role, i)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(role)
}

func (h Handler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	var roles []Roles
	h.conn.GetConnection().Find(&roles)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(roles)
}

func MustNewHandlerRole() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
