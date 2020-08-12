package school_service

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

func (h Handler) GetAllSchools(w http.ResponseWriter, r *http.Request) {
	var schools []School
	h.conn.GetConnection().Find(&schools)

	_ = json.NewEncoder(w).Encode(schools)
}

func (h Handler) AddSchool(w http.ResponseWriter, r *http.Request) {
	var school School
	//@todo: cover error
	_ = json.NewDecoder(r.Body).Decode(&school)

	h.conn.GetConnection().Preload("Owner").Save(&school).Find(&school)

	_ = json.NewEncoder(w).Encode(school)
}

func (h Handler) GetSchoolById(w http.ResponseWriter, r *http.Request) {
	var school School

	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	h.conn.GetConnection().Preload("Owner").Find(&school, i)

	_ = json.NewEncoder(w).Encode(school)
}

func (h Handler) UpdateSchoolById(w http.ResponseWriter, r *http.Request) {
	var school School

	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	//@todo: cover error
	var updates = make(map[string]interface{})
	_ = json.NewDecoder(r.Body).Decode(&updates)

	err = h.conn.GetConnection().First(&school, i).Error
	if err != nil {
		h.log.Error("error on update school", err, nil)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = h.conn.GetConnection().Model(&school).Updates(updates).Error
	if err != nil {
		h.log.Error("error on update school", err, nil)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_ = json.NewEncoder(w).Encode(school)
}

func (h Handler) DeleteSchoolById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	h.conn.GetConnection().Delete(School{}, i)

	_ = json.NewEncoder(w).Encode("School deleted")
}

func MustNewHandlerSchool() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
