package direction_service

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

func (h Handler) DisbandDirectionById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")
		return
	}

	h.conn.GetConnection().Where(Direction{}, i).Delete(&Direction{})

	_ = json.NewEncoder(w).Encode("direction deleted")
}

func (h Handler) AddDirection(w http.ResponseWriter, r *http.Request) {
	var d Direction

	// @todo: work with error
	_ = json.NewDecoder(r.Body).Decode(&d)

	err := h.conn.GetConnection().Save(&d).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_ = json.NewEncoder(w).Encode(d)
}

func (h Handler) UpdateDirectionById(w http.ResponseWriter, r *http.Request) {
	var newDirection Direction
	var oldDirection Direction
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	// @todo: work with error
	_ = json.NewDecoder(r.Body).Decode(&newDirection)

	h.conn.GetConnection().Where(&oldDirection, i).Update(newDirection)

	_ = json.NewEncoder(w).Encode(oldDirection)
}

func (h Handler) GetDirectionById(w http.ResponseWriter, r *http.Request) {
	var direction Direction
	i, _ := strconv.Atoi(mux.Vars(r)["id"])

	h.conn.GetConnection().First(&direction, i)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(direction)
}

func (h Handler) GetAllDirections(w http.ResponseWriter, _ *http.Request) {
	var directions []Direction

	h.conn.GetConnection().Find(&directions)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(directions)
}

func MustNewHandlerDirection() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
