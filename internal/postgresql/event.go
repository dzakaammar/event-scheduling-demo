package postgresql

import (
	"context"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (e *EventRepository) Store(ctx context.Context, event *internal.Event) error {
	return nil
}

func (e *EventRepository) DeleteByID(ctx context.Context, id string) error {
	return nil
}

func (e *EventRepository) Update(ctx context.Context, event *internal.Event) error {
	return nil
}

func (e *EventRepository) FindByID(ctx context.Context, id string) (*internal.Event, error) {
	return nil, nil
}
