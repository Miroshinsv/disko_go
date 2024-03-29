package user_service

import (
	"encoding/json"
	rolesModel "github.com/Miroshinsv/disko_go/internal/role-service"
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

func (h Handler) DisbandUserById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	h.conn.GetConnection().Where(Users{}, i).Delete(&Users{})

	_ = json.NewEncoder(w).Encode("user disband")
}

func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var d Users
	var roles rolesModel.Roles

	//@todo: cover error
	_ = json.NewDecoder(r.Body).Decode(&d)

	res := h.conn.GetConnection().Find(&roles).Where(rolesModel.Roles{Admin: false})

	if res == nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err := h.conn.GetConnection().Save(&d).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)

		return
	}

	_ = json.NewEncoder(w).Encode(d)
}

func (h Handler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	var (
		newUser Users
		oldUser Users
	)

	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	//@todo: cover error
	_ = json.NewDecoder(r.Body).Decode(&newUser)

	h.conn.GetConnection().Where(&oldUser, i).Update(newUser)

	_ = json.NewEncoder(w).Encode(oldUser)
}

func (h Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	var user Users

	//@todo: cover error
	i, _ := strconv.Atoi(mux.Vars(r)["id"])

	h.conn.GetConnection().Preload("Roles").First(&user, i)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
}

func (h Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []Users
	h.conn.GetConnection().Preload("Roles").Find(&users)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}

func MustNewHandlerUser() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
