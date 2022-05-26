package scheduling

import (
	"context"
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal/core"
)

type EventService struct {
	eventRepo core.EventRepository
}

func NewEventService(eventRepo core.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

func (e *EventService) CreateEvent(ctx context.Context, req *core.CreateEventRequest) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	err = e.eventRepo.Store(ctx, req.Event)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventService) DeleteEventByID(ctx context.Context, req *core.DeleteEventByIDRequest) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	err = e.eventRepo.DeleteByID(ctx, req.EventID)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventService) UpdateEvent(ctx context.Context, req *core.UpdateEventRequest) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	now := time.Now()
	req.Event.UpdatedAt = &now

	err = e.eventRepo.Update(ctx, req.Event)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventService) FindEventByID(ctx context.Context, req *core.FindEventByIDRequest) (*core.Event, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	event, err := e.eventRepo.FindByID(ctx, req.EventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}