package event_service

import (
	"encoding/json"
	"fmt"
	"github.com/Miroshinsv/disko_go/internal/event-service/models"
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Handler struct {
	log          loggerService.ILogger
	conn         dbConnector.IConnector
	eventService *Service
}

var (
	invalidEventId  = errors.New("Invalid event ID")
	cantUpdateEvent = errors.New("Can't update event")
	eventUpdated    = "Event update"
)

func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})
}

func (h Handler) DeleteEventById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(invalidEventId.Error())
		return
	}

	h.conn.GetConnection().Where(models.Events{}, i).Delete(&models.Events{}, i)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(invalidEventId.Error())
}

func (h Handler) DeactivateEventById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(invalidEventId.Error())

		return
	}

	var events models.Events
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
		_ = json.NewEncoder(w).Encode(invalidEventId.Error())

		return
	}

	var events models.Events
	h.conn.GetConnection().Preload("Type").Find(&events, i).Updates(models.Events{IsActive: true})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) GetEventById(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(invalidEventId.Error())

		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var events models.Events
	h.conn.GetConnection().
		Preload("Type").
		Preload("Polls").
		Find(&events, i)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) GetAllEvents(w http.ResponseWriter, _ *http.Request) {
	var events []models.Events
	h.conn.GetConnection().
		Preload("Type").
		Preload("Polls").
		Find(&events)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) AddEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Events
	rErr := json.NewDecoder(r.Body).Decode(&event)
	if rErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid event: " + rErr.Error())
		return
	}

	event.OwnerId = r.Context().Value("user").(*userService.Users).ID

	err := h.conn.GetConnection().Save(&event)
	if err.Error != nil {
		_ = json.NewEncoder(w).Encode(cantUpdateEvent.Error())
	} else {
		h.conn.GetConnection().Preload("Type").Find(&event, event.ID)
		_ = json.NewEncoder(w).Encode(event)
	}
}

func (h Handler) UpdateEventById(w http.ResponseWriter, r *http.Request) {
	var nEvent models.Events

	i, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(invalidEventId.Error())

		return
	}
	//@todo: cover error
	_ = json.NewDecoder(r.Body).Decode(&nEvent)
	err = h.eventService.Update(i, r.Context().Value("user").(*userService.Users).ID, &nEvent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(eventUpdated)
}

func (h Handler) GetEventsType(w http.ResponseWriter, r *http.Request) {
	var nEventTypes []models.EventsType

	err := h.eventService.FindEventsType(&nEventTypes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	_ = json.NewEncoder(w).Encode(nEventTypes)
}

func MustNewHandlerEvent() *Handler {
	db, _ := dbConnector.GetDBConnection()
	log := loggerService.GetLogger()

	return &Handler{
		log:          loggerService.GetLogger(),
		conn:         db,
		eventService: MustNewEventService(log, db),
	}
}
