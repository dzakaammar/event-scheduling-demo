package postgresql

import (
	"context"
	"database/sql"

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
	trx, err := e.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer trx.Rollback()

	res, err := trx.Exec(`INSERT INTO public.event (id, title, description, timezone, created_by, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`, event.ID, event.Title, event.Description, event.Timezone, event.CreatedBy, event.CreatedAt)
	if err != nil {
		return err
	}

	return trx.Commit()
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
