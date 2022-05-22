package postgresql

import (
	"context"

	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Instrumentation struct {
	next   core.EventRepository
	tracer trace.Tracer
}

func NewInstrumentation(next core.EventRepository) *Instrumentation {
	return &Instrumentation{
		next:   next,
		tracer: otel.Tracer("event-repository"),
	}
}

func (i *Instrumentation) Store(ctx context.Context, event *core.Event) error {
	var err error
	ctx, span := i.tracer.Start(ctx, "store")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = i.next.Store(ctx, event)
	return err
}

func (i *Instrumentation) DeleteByID(ctx context.Context, id string) error {
	var err error
	ctx, span := i.tracer.Start(ctx, "delete-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = i.next.DeleteByID(ctx, id)
	return err
}

func (i *Instrumentation) Update(ctx context.Context, event *core.Event) error {
	var err error
	ctx, span := i.tracer.Start(ctx, "update")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	err = i.next.Update(ctx, event)
	return err
}

func (i *Instrumentation) FindByID(ctx context.Context, id string) (*core.Event, error) {
	var err error
	ctx, span := i.tracer.Start(ctx, "find-by-id")
	defer func() {
		span.RecordError(err)
		span.End()
	}()

	event, err := i.next.FindByID(ctx, id)
	return event, err
}
