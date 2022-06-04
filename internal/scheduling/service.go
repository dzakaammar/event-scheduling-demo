package scheduling

import (
	"context"
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal/core"
)

type Service struct {
	eventRepo core.EventRepository
}

func NewService(eventRepo core.EventRepository) *Service {
	return &Service{
		eventRepo: eventRepo,
	}
}

func (e *Service) CreateEvent(ctx context.Context, req *core.CreateEventRequest) error {
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

func (e *Service) DeleteEventByID(ctx context.Context, req *core.DeleteEventByIDRequest) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	_, err = e.eventRepo.FindByID(ctx, req.EventID)
	if err != nil {
		return err
	}

	err = e.eventRepo.DeleteByID(ctx, req.EventID)
	if err != nil {
		return err
	}
	return nil
}

func (e *Service) UpdateEvent(ctx context.Context, req *core.UpdateEventRequest) error {
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

func (e *Service) FindEventByID(ctx context.Context, req *core.FindEventByIDRequest) (*core.Event, error) {
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
