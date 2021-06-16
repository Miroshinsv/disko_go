package poll_service

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"

	eventModel "github.com/Miroshinsv/disko_go/internal/event-service/models"
	"github.com/Miroshinsv/disko_go/internal/poll-service/models"
	userService "github.com/Miroshinsv/disko_go/internal/user-service"
	dbConnector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
)

const (
	timeFormat = "2006-01-02"
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
		event    = &eventModel.Events{}
	)

	s.conn.GetConnection().Where(
		fmt.Sprintf("event_id=%d", p.EventId),
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

func (s Service) Vote(poll *models.Poll, user *userService.Users) error {
	var vote = &models.Vote{}
	s.conn.GetConnection().Where(fmt.Sprintf("poll_id=%d AND user_id=%d", poll.ID, user.ID)).Find(vote)
	if vote.ID != 0 {
		return errorAlreadyVoted
	}

	return s.conn.GetConnection().Create(&models.Vote{
		PollId:    int(poll.ID),
		UserId:    int(user.ID),
		CreatedAt: time.Now(),
	}).Error
}

func (s Service) ShowResults(poll *models.Poll, user *userService.Users) ([]models.Vote, error) {
	if time.Now().Before(poll.DueDate) && poll.IsHidden {
		err := s.conn.GetConnection().Where(fmt.Sprintf("poll_id=%d AND user_id=%d", poll.ID, user.ID)).Find(&models.Vote{}).Error
		if gorm.IsRecordNotFoundError(err) {
			return []models.Vote{}, errorVoteNeeded
		}
	}

	var votes []models.Vote
	db := s.conn.GetConnection().Where(fmt.Sprintf("poll_id=%d", poll.ID)).Preload("User").Find(&votes)
	if db.Error != nil {
		return []models.Vote{}, db.Error
	}

	return votes, nil
}

func (s Service) ShowVotesCount(poll *models.Poll) (int, error) {
	var votes []models.Vote
	db := s.conn.GetConnection().Where(fmt.Sprintf("poll_id=%d", poll.ID)).Preload("User").Find(&votes)
	if db.Error != nil {
		return 0, db.Error
	}

	return len(votes), nil
}

func (s Service) ScheduleAutoPolls(events []eventModel.Events, dt time.Time) error {
	var (
		eventIds = make([]uint, 0)
		plEIds   = make(map[int]models.Poll, 0)
		polls    []models.Poll
		ds       = dt.Format(timeFormat)
	)

	for _, v := range events {
		eventIds = append(eventIds, v.ID)
	}

	s.conn.GetConnection().
		Where("event_id IN (?)", eventIds).
		Where("TO_CHAR(due_date, 'YYYY-MM-DD') = ?", ds).
		Find(&polls)

	for _, v := range polls {
		plEIds[v.EventId] = v
	}

	for _, i := range eventIds {
		if _, isExists := plEIds[int(i)]; isExists {
			continue
		}

		dbPoll := &models.Poll{
			Model:    gorm.Model{},
			EventId:  int(i),
			IsHidden: true,
			DueDate:  dt,
		}

		db := s.conn.GetConnection().Create(&dbPoll)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
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
