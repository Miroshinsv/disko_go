package schedule_service

import (
	"fmt"
	eventService "github.com/Miroshinsv/disko_go/internal/event-service"
	db_connector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	"strings"
	"time"
)

var self *Service = nil

type Service struct{}

func (s Service) LoadEventsForDate(d time.Time) ([]eventService.Events, error) {
	var (
		result = make([]eventService.Events, 0)
	)

	conn, err := db_connector.GetDBConnection()
	if err != nil {
		return result, err
	}

	conn.GetConnection().Preload("Type").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Where(fmt.Sprintf("'%s' = ANY(events.days) AND events.is_active = true", names[d.Weekday()])).
		Find(&result)

	return result, nil
}

func (s Service) LoadEventsForPeriod(from time.Time, to time.Time) (map[string][]eventService.Events, error) {
	var (
		result  = make(map[string][]eventService.Events, 0)
		events  []eventService.Events
		tEvents []eventService.Events
	)

	conn, err := db_connector.GetDBConnection()
	if err != nil {
		return result, err
	}

	conn.GetConnection().Preload("Type").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Where("events.is_active = true AND events_types.is_repeatable = true").
		Find(&events)

	conn.GetConnection().Preload("Type").
		Joins("LEFT JOIN events_types ON events.type_id = events_types.id").
		Where(
			fmt.Sprintf(
				"events.created_at BETWEEN '%s' AND '%s' AND events.is_active = true AND events_types.is_repeatable = false",
				from.Format(timeFormat), to.Format(timeFormat),
			),
		).
		Find(&tEvents)

	ev := append(events, tEvents...)
	days := to.Sub(from).Hours() / 24
	for i := 0; i < (int(days) + 1); i++ {
		curDate := from.AddDate(0, 0, i)
		result[from.AddDate(0, 0, i).Format(timeFormat)] = s.findEventsForDate(curDate, ev)
	}

	return result, nil
}

func (s Service) findEventsForDate(d time.Time, events []eventService.Events) []eventService.Events {
	var res = make([]eventService.Events, 0)

	weekDay := strings.ToLower(d.Weekday().String())
	for _, v := range events {
		if strings.Contains(v.Days, weekDay) {
			res = append(res, v)
		}
	}

	return res
}

func GetScheduleService() *Service {
	if self == nil {
		self = &Service{}
	}

	return self
}
