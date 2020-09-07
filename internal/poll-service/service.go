package poll_service

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	event_service "github.com/Miroshinsv/disko_go/internal/event-service"
	"github.com/Miroshinsv/disko_go/internal/poll-service/models"
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
)

var (
	self *Service = nil

	errorExistingPoll = errors.New("same poll already exists for this event")
	errorUnknownEvent = errors.New("unknown event ID")
	errorUnknownVoice = errors.New("unknown voice type")
	errorAlreadyVoted = errors.New("already voted")
	errorVoteNeeded   = errors.New("vote needed")
)

type Service struct {
	log  loggerService.ILogger
	conn dbConnector.IConnector
}

func (s Service) Create(p models.Income) (*models.Poll, error) {
	var (
		existing = &models.Poll{}
		event    = &event_service.Events{}
	)

	s.conn.GetConnection().Where(
		fmt.Sprintf("name LIKE '%s' AND event_id=%d", p.Name, p.EventId),
	).Find(existing)
	if existing.ID != 0 {
		return &models.Poll{}, errorExistingPoll
	}

	s.conn.GetConnection().Where(fmt.Sprintf("id=%d", p.EventId)).Find(event)
	if event.ID == 0 {
		return &models.Poll{}, errorUnknownEvent
	}

	dbPoll := &models.Poll{
		Model:    gorm.Model{},
		Name:     p.Name,
		Subject:  p.Subject,
		EventId:  int(event.ID),
		IsHidden: p.IsHidden,
		DueDate:  time.Unix(p.DueDate, 0),
	}

	db := s.conn.GetConnection().Create(&dbPoll)
	if db.Error != nil {
		return dbPoll, db.Error
	}

	return dbPoll, nil
}

func (s Service) Vote(voice int, poll *models.Poll, user *userService.Users) error {
	if !s.isVoiceValid(voice) {
		return errorUnknownVoice
	}

	var vote = &models.Vote{}
	s.conn.GetConnection().Where(fmt.Sprintf("poll_id=%d AND user_id=%d", poll.ID, user.ID)).Find(vote)
	if vote.ID != 0 {
		return errorAlreadyVoted
	}

	return s.conn.GetConnection().Create(&models.Vote{
		PollId:    int(poll.ID),
		UserId:    int(user.ID),
		Voice:     voice,
		CreatedAt: time.Now(),
	}).Error
}

func (s Service) ShowResults(poll *models.Poll, user *userService.Users) (map[int][]models.Vote, error) {
	if time.Now().Before(poll.DueDate) && poll.IsHidden {
		var vote = &models.Vote{}
		s.conn.GetConnection().Where(fmt.Sprintf("poll_id=%d AND user_id=%d", poll.ID, user.ID)).Find(vote)
		if vote.ID == 0 {
			return map[int][]models.Vote{}, errorVoteNeeded
		}
	}

	var votes []models.Vote
	db := s.conn.GetConnection().Where(fmt.Sprintf("poll_id=%d", poll.ID)).Preload("User").Find(&votes)
	if db.Error != nil {
		return map[int][]models.Vote{}, db.Error
	}

	var result = make(map[int][]models.Vote, 3)
	for _, v := range votes {
		result[v.Voice] = append(result[v.Voice], v)
	}

	return result, nil
}

func (s Service) isVoiceValid(voice int) bool {
	for _, v := range []int{models.PositiveVoice, models.NegativeVoice, models.NeutralVoice} {
		if v == voice {
			return true
		}
	}

	return false
}

func MustNewPollService(log loggerService.ILogger, conn dbConnector.IConnector) *Service {
	if self == nil {
		self = &Service{
			log:  log,
			conn: conn,
		}
	}

	return self
}

func GetPollService() *Service {
	return self
}
