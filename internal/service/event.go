package service

import (
	"MyStonks-go/internal/models"
	"MyStonks-go/internal/store"
)

type EventSrv struct {
	es store.EventStore
}

func NewEventSrv(store store.EventStore) *EventSrv {
	return &EventSrv{es: store}
}

func (s *EventSrv) GetEvents() ([]models.Events, error) {
	return s.es.GetEvents()
}
