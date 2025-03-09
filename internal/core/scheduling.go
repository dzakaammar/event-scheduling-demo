package core

import (
	"context"

	"github.com/dzakaammar/event-scheduling-example/internal"
)

type CreateEventRequest struct {
	ActorID string
	Event   *Event
}

func (c *CreateEventRequest) Validate() error {
	if c.ActorID == "" {
		return internal.WrapErr(internal.ErrValidationFailed, "invalid actor id")
	}

	if c.Event == nil {
		return internal.WrapErr(internal.ErrValidationFailed, "invalid event")
	}

	if len(c.Event.Schedules) <= 0 {
		return internal.WrapErr(internal.ErrValidationFailed, "no schedules provided for the event")
	}

	return c.Event.Validate()
}

type DeleteEventByIDRequest struct {
	ActorID string
	EventID string
}

func (d *DeleteEventByIDRequest) Validate() error {
	if d.ActorID == "" {
		return internal.WrapErr(internal.ErrValidationFailed, "invalid actor id")
	}

	if d.EventID == "" {
		return internal.WrapErr(internal.ErrValidationFailed, "invalid event id")
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
		return internal.WrapErr(internal.ErrValidationFailed, "invalid actor id")
	}

	if u.Event == nil {
		return internal.WrapErr(internal.ErrValidationFailed, "invalid event")
	}

	return u.Event.Validate()
}

type FindEventByIDRequest struct {
	EventID string
}

func (f *FindEventByIDRequest) Validate() error {
	if f.EventID == "" {
		return internal.WrapErr(internal.ErrValidationFailed, "invalid event id")
	}

	return nil
}

//go:generate go tool -modfile=../../go.tool.mod mockgen -destination=../mock/mock_scheduling_service.go -package=mock github.com/dzakaammar/event-scheduling-example/internal/core SchedulingService
type SchedulingService interface {
	CreateEvent(ctx context.Context, req *CreateEventRequest) error
	DeleteEventByID(ctx context.Context, req *DeleteEventByIDRequest) error
	UpdateEvent(ctx context.Context, req *UpdateEventRequest) error
	FindEventByID(ctx context.Context, req *FindEventByIDRequest) (*Event, error)
}
