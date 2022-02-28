package internal

import (
	"context"
)

type CreateEventRequest struct {
	ActorID string
	Event   *Event
}

func (c *CreateEventRequest) Validate() error {
	if c.ActorID == "" {
		return WrapErr(ErrValidationFailed, "invalid actor id")
	}

	if c.Event == nil {
		return WrapErr(ErrValidationFailed, "invalid event")
	}

	return c.Event.Validate()
}

type DeleteEventByIDRequest struct {
	ActorID string
	EventID string
}

func (d *DeleteEventByIDRequest) Validate() error {
	if d.ActorID == "" {
		return WrapErr(ErrValidationFailed, "invalid actor id")
	}

	if d.EventID == "" {
		return WrapErr(ErrValidationFailed, "invalid event id")
	}

	return nil
}

type UpdateEventRequest struct {
	ID      string
	ActorID string
	Event   *Event
}

func (u *UpdateEventRequest) Validate() error {
	if u.ActorID == "" {
		return WrapErr(ErrValidationFailed, "invalid actor id")
	}

	if u.Event == nil {
		return WrapErr(ErrValidationFailed, "invalid event")
	}

	return u.Event.Validate()
}

type FindEventByIDRequest struct {
	EventID string
}

func (f *FindEventByIDRequest) Validate() error {
	if f.EventID == "" {
		return WrapErr(ErrValidationFailed, "invalid event id")
	}

	return nil
}

//go:generate mockgen -destination=mock/mock_event_service.go -package=mock github.com/dzakaammar/event-scheduling-example/internal EventService
type EventService interface {
	CreateEvent(ctx context.Context, req *CreateEventRequest) error
	DeleteEventByID(ctx context.Context, req *DeleteEventByIDRequest) error
	UpdateEvent(ctx context.Context, req *UpdateEventRequest) error
	FindEventByID(ctx context.Context, req *FindEventByIDRequest) (*Event, error)
}
