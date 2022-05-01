package instrumentation

import (
	"context"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type EventRepository struct {
	next   internal.EventRepository
	tracer trace.Tracer
}

func NewEventRepository(next internal.EventRepository) *EventRepository {
	return &EventRepository{
		next:   next,
		tracer: otel.Tracer("event-repository"),
	}
}

func (e *EventRepository) Store(ctx context.Context, event *internal.Event) error {
	var err error
	ctx, span := e.tracer.Start(ctx, "store")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = e.next.Store(ctx, event)
	return err
}

func (e *EventRepository) DeleteByID(ctx context.Context, id string) error {
	var err error
	ctx, span := e.tracer.Start(ctx, "delete-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = e.next.DeleteByID(ctx, id)
	return err
}

func (e *EventRepository) Update(ctx context.Context, event *internal.Event) error {
	var err error
	ctx, span := e.tracer.Start(ctx, "update")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = e.next.Update(ctx, event)
	return err
}

func (e *EventRepository) FindByID(ctx context.Context, id string) (*internal.Event, error) {
	var err error
	ctx, span := e.tracer.Start(ctx, "find-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	event, err := e.next.FindByID(ctx, id)
	return event, err
}
