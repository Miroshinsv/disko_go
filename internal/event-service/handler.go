package event_service

import (
	"encoding/json"
	"fmt"
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

func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func (h Handler) DeleteEventById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")
		return
	}

	h.conn.GetConnection().Where(Events{}, i).Delete(&Events{}, i)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Event disband")
}

func (h Handler) DeactivateEventById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	var events Events
	res := h.conn.GetConnection().Preload("Type").Find(&events, i)
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(res.Error.Error())

		return
	}

	h.conn.GetConnection().Model(&events).Updates(map[string]interface{}{
		"is_active": false,
	})
	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(res.Error.Error())

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) ActivateEventById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	var events Events
	h.conn.GetConnection().Preload("Type").Find(&events, i).Updates(Events{IsActive: true})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) GetEventById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var events Events
	h.conn.GetConnection().
		Preload("Type").
		Preload("Polls").
		Find(&events, i)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) GetAllEvents(w http.ResponseWriter, _ *http.Request) {
	var events []Events
	h.conn.GetConnection().
		Preload("Type").
		Preload("Polls").
		Find(&events)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) AddEvent(w http.ResponseWriter, r *http.Request) {
	var event Events

	rErr := json.NewDecoder(r.Body).Decode(&event)
	if rErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid event: " + rErr.Error())

		return
	}

	err := h.conn.GetConnection().Save(&event)
	if err.Error != nil {
		_ = json.NewEncoder(w).Encode("Can't add event")
	} else {
		h.conn.GetConnection().Preload("Type").Find(&event, event.ID)
		_ = json.NewEncoder(w).Encode(event)
	}
}

func (h Handler) UpdateEventById(w http.ResponseWriter, r *http.Request) {
	var (
		nEvent Events
		event  Events
	)

	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid id")

		return
	}

	//@todo: cover error
	_ = json.NewDecoder(r.Body).Decode(&nEvent)

	h.conn.GetConnection().Find(&event, i).Updates(nEvent)
	_ = json.NewEncoder(w).Encode(event)
}

func MustNewHandlerEvent() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
