package models

import (
	"errors"
	"time"
)

var (
	errorEmptyName    = errors.New("empty name for poll")
	errorEmptySubject = errors.New("empty subject for poll")
	errorDueDate      = errors.New("incorrect due date for poll")
	errorEventId      = errors.New("incorrect event ID for poll")
)

type Income struct {
	Name     string
	EventId  int `json:"event_id"`
	Subject  string
	DueDate  int64 `json:"due_date"`
	IsHidden bool
}

func (s Income) Validate() error {
	if s.Name == "" {
		return errorEmptyName
	}

	if s.EventId == 0 {
		return errorEventId
	}

	if s.Subject == "" {
		return errorEmptySubject
	}

	if s.DueDate == 0 {
		return errorDueDate
	}

	if time.Now().Unix() >= s.DueDate {
		return errorDueDate
	}

	return nil
}
