package instrumentation

import (
	"context"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type EventService struct {
	next   internal.EventService
	tracer trace.Tracer
}

func NewEventService(next internal.EventService) *EventService {
	return &EventService{
		next:   next,
		tracer: otel.Tracer("event-service"),
	}
}

func (e *EventService) CreateEvent(ctx context.Context, req *internal.CreateEventRequest) error {
	var err error
	ctx, span := e.tracer.Start(ctx, "create-event")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = e.next.CreateEvent(ctx, req)
	return err
}

func (e *EventService) DeleteEventByID(ctx context.Context, req *internal.DeleteEventByIDRequest) error {
	var err error
	ctx, span := e.tracer.Start(ctx, "delete-event-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = e.next.DeleteEventByID(ctx, req)
	return err
}

func (e *EventService) UpdateEvent(ctx context.Context, req *internal.UpdateEventRequest) error {
	var err error
	ctx, span := e.tracer.Start(ctx, "update-event")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = e.next.UpdateEvent(ctx, req)
	return err
}

func (e *EventService) FindEventByID(ctx context.Context, req *internal.FindEventByIDRequest) (*internal.Event, error) {
	var err error
	ctx, span := e.tracer.Start(ctx, "find-event-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	event, err := e.next.FindEventByID(ctx, req)
	return event, err
}
