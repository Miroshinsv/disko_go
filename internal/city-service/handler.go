package city_service

import (
	"encoding/json"
	"github.com/Miroshinsv/disko_go/internal/city-service/models"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	log  loggerService.ILogger
	conn dbConnector.IConnector
}

func (h Handler) DisbandCityById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")
		return
	}

	h.conn.GetConnection().Where(models.City{}, i).Delete(&models.City{})

	_ = json.NewEncoder(w).Encode("city deleted")
}

func (h Handler) AddCity(w http.ResponseWriter, r *http.Request) {
	var d models.City

	// @todo: work with error
	_ = json.NewDecoder(r.Body).Decode(&d)

	d.CityName = strings.Title(strings.ToLower(d.CityName))

	err := h.conn.GetConnection().Save(&d).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	_ = json.NewEncoder(w).Encode(d)
}

func (h Handler) UpdateCityById(w http.ResponseWriter, r *http.Request) {
	var newCity models.City
	var oldCity models.City
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	// @todo: work with error
	_ = json.NewDecoder(r.Body).Decode(&newCity)

	h.conn.GetConnection().Where(&oldCity, i).Update(newCity)

	_ = json.NewEncoder(w).Encode(oldCity)
}

func (h Handler) GetCityById(w http.ResponseWriter, r *http.Request) {
	var city models.City
	i, _ := strconv.Atoi(mux.Vars(r)["id"])

	h.conn.GetConnection().First(&city, i)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(city)
}

func (h Handler) GetAllCities(w http.ResponseWriter, _ *http.Request) {
	var events []models.City
	h.conn.GetConnection().
		Find(&events)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func MustNewHandlerCities() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
