package scheduling

import (
	"context"

	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Instrumentation struct {
	next   core.SchedulingService
	tracer trace.Tracer
}

func NewInstrumentation(next core.SchedulingService) *Instrumentation {
	return &Instrumentation{
		next:   next,
		tracer: otel.Tracer("event-service"),
	}
}

func (i *Instrumentation) CreateEvent(ctx context.Context, req *core.CreateEventRequest) error {
	var err error
	ctx, span := i.tracer.Start(ctx, "create-event")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = i.next.CreateEvent(ctx, req)
	return err
}

func (i *Instrumentation) DeleteEventByID(ctx context.Context, req *core.DeleteEventByIDRequest) error {
	var err error
	ctx, span := i.tracer.Start(ctx, "delete-event-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = i.next.DeleteEventByID(ctx, req)
	return err
}

func (i *Instrumentation) UpdateEvent(ctx context.Context, req *core.UpdateEventRequest) error {
	var err error
	ctx, span := i.tracer.Start(ctx, "update-event")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = i.next.UpdateEvent(ctx, req)
	return err
}

func (i *Instrumentation) FindEventByID(ctx context.Context, req *core.FindEventByIDRequest) (*core.Event, error) {
	var err error
	ctx, span := i.tracer.Start(ctx, "find-event-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	event, err := i.next.FindEventByID(ctx, req)
	return event, err
}
