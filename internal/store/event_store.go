package store

import (
	"MyStonks-go/internal/models"
	"gorm.io/gorm"
)

type EventStore interface {
	GetEvents() ([]models.Events, error)
}

type eventStore struct {
	db *gorm.DB
}

func NewEventStore(db *gorm.DB) EventStore {
	return &eventStore{db: db}
}

func (e *eventStore) GetEvents() ([]models.Events, error) {
	var events []models.Events
	if err := e.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
