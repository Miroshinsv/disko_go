package event_service

import (
	"fmt"
	"github.com/Miroshinsv/disko_go/internal/event-service/models"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	"github.com/pkg/errors"
)

type Service struct {
	log  loggerService.ILogger
	conn dbConnector.IConnector
}

var (
	self             *Service = nil
	errorUpdateEvent          = errors.New("Can't update event")
)

func (s Service) Update(eId int, uId uint, event *models.Events) error {
	err := s.conn.GetConnection().Model(event).Where(fmt.Sprintf("id=%d AND owner_id=%d", eId, uId)).Update(event)
	if err.RowsAffected == 0 {
		return errorUpdateEvent
	}
	return nil
}

func MustNewEventService(log loggerService.ILogger, conn dbConnector.IConnector) *Service {
	if self == nil {
		self = &Service{
			log:  log,
			conn: conn,
		}
	}

	return self
}
