package schedule_service

import (
	"encoding/json"
	"fmt"
	eventService "github.com/Miroshinsv/disko_go/internal/event-service"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"net/http"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

var names = [...]string{
	"sunday",
	"monday",
	"tuesday",
	"wednesday",
	"thursday",
	"friday",
	"saturday",
}

type Handler struct {
	log  loggerService.ILogger
	conn dbConnector.IConnector
}

func (h Handler) LoadAllEvents(w http.ResponseWriter, _ *http.Request) {
	var events []eventService.Events
	h.conn.GetConnection().Preload("Type").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Find(
			&events,
			fmt.Sprintf("events.is_active = true"),
		)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(events)
}

func (h Handler) LoadEventsForToday(w http.ResponseWriter, _ *http.Request) {
	var events []eventService.Events
	h.conn.GetConnection().Preload("Type").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Find(
			&events,
			fmt.Sprintf(
				"events.days = '{%s}' AND events.is_active = true",
				names[time.Now().Weekday()],
			),
		)

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

	var events []eventService.Events
	h.conn.GetConnection().Preload("Type").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Find(&events,
			fmt.Sprintf(
				"events.created_at <= '%s' AND events.is_active = true AND events_types.day_of_week != dayofweek('%d') AND events_types.is_repeatable = true",
				dateTo.Format(timeFormat),
				time.Now().Weekday(),
			),
		)

	var tEvents []eventService.Events
	h.conn.GetConnection().Preload("Type").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Find(
			&tEvents,
			fmt.Sprintf(
				"events.created_at <= NOW() AND events.is_active = true AND events_types.day_of_week = dayofweek('%d') AND is_repeatable = false",
				time.Now().Weekday(),
			),
		)

	ev := append(events, tEvents...)

	var result = make(map[string][]eventService.Events)

	days := dateTo.Sub(dateFrom).Hours() / 24
	for i := 0; i < int(days); i++ {
		curDate := dateFrom.AddDate(0, 0, i)
		result[dateFrom.AddDate(0, 0, i).Format(timeFormat)] = h.findEventsForDate(curDate, ev)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func (h Handler) findEventsForDate(d time.Time, events []eventService.Events) []eventService.Events {
	var res = make([]eventService.Events, 0)

	//weekDay := int(d.Weekday())
	//for _, v := range events {
	//	if v.Type.DayOfWeek == weekDay {
	//		res = append(res, v)
	//	}
	//}

	return res
}

func MustNewHandlerSchedule() *Handler {
	db, _ := dbConnector.GetDBConnection()

	return &Handler{
		log:  loggerService.GetLogger(),
		conn: db,
	}
}
