package schedule_service

import (
	"encoding/json"
	"github.com/Miroshinsv/disko_go/internal/event-service/models"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"net/http"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

type Handler struct {
	log  loggerService.ILogger
	conn dbConnector.IConnector
}

func (h Handler) LoadAllEvents(w http.ResponseWriter, r *http.Request) {
	var (
		events []models.Events
	)

	stm := h.conn.GetConnection().Preload("Type").
		Preload("Polls").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Where("events.is_active = true").
		Preload("City").
		Preload("Dj")

	city := r.URL.Query().Get("city")

	if city == "" {
		stm.
			Find(
				&events,
			)

	} else {
		stm.
			Where("city_name = ?", strings.Title(strings.ToLower(city))).
			Find(
				&events,
			)

	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) LoadEventsForToday(w http.ResponseWriter, _ *http.Request) {
	scheduleService := GetScheduleService()
	events, err := scheduleService.LoadEventsForDate(time.Now())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) LoadEventsForPeriod(w http.ResponseWriter, r *http.Request) {
	dateFrom, err := time.Parse(timeFormat, r.URL.Query().Get("from"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
	}

	dateTo, err := time.Parse(timeFormat, r.URL.Query().Get("to"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
	}

	scheduleService := GetScheduleService()
	events, err := scheduleService.LoadEventsForPeriod(dateFrom, dateTo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func MustNewHandlerSchedule() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
