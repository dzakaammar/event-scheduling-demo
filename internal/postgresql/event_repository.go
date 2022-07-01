package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
)

type EventRepository struct {
	dbConn *sqlx.DB
}

func NewEventRepository(dbConn *sqlx.DB) *EventRepository {
	return &EventRepository{
		dbConn: dbConn,
	}
}

func (e *EventRepository) Store(ctx context.Context, event *core.Event) error {
	log := logger.WithFields(logger.Fields{
		"event": event,
	})

	tx, err := e.dbConn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != sql.ErrTxDone {
			log.Error(rollbackErr)
		}
	}()

	createEventSql := "INSERT INTO event (id, title, description, timezone, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = tx.ExecContext(ctx, createEventSql, event.ID, event.Title, event.Description, event.Timezone, event.CreatedBy, time.Now(), time.Now())
	if err != nil {
		log.Error(err)
		return err
	}

	createScheduleSql := "INSERT INTO schedule (id, event_id, start_time, duration, is_full_day, recurring_interval, recurring_type) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	for _, schedule := range event.Schedules {
		_, err = tx.ExecContext(ctx, createScheduleSql, schedule.ID, event.ID, schedule.StartTime, schedule.DurationInMinutes, schedule.IsFullDay, schedule.RecurringInterval, schedule.RecurringType)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	createInvitationSql := "INSERT INTO invitation (id, event_id, user_id, token, status) VALUES ($1, $2, $3, $4, $5)"
	for _, invitation := range event.Invitations {
		_, err = tx.ExecContext(ctx, createInvitationSql, invitation.ID, event.ID, invitation.UserID, invitation.Token, invitation.Status)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return tx.Commit()
}

func (e *EventRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := e.dbConn.ExecContext(ctx, "DELETE FROM event WHERE id = $1", id)
	if err != nil {
		logger.WithField("id", id).Error(err)
		return err
	}
	return nil
}

func (e *EventRepository) Update(ctx context.Context, event *core.Event) error {
	log := logger.WithFields(logger.Fields{
		"event": event,
	})

	tx, err := e.dbConn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != sql.ErrTxDone {
			log.Error(rollbackErr)
		}
	}()

	updateSql := `UPDATE event SET title = $1, description = $2, timezone = $3, updated_at = $4`
	_, err = tx.ExecContext(ctx, updateSql, event.Title, event.Description, event.Timezone, time.Now())
	if err != nil {
		log.Error(err)
		return nil
	}

	upsertScheduleSql := `INSERT INTO schedule (id, event_id, start_time, duration, is_full_day, recurring_interval, recurring_type) VALUES ($1, $2, $3, $4, $5, $6, $7)
						ON CONFLICT (id, event_id)
						DO UPDATE SET start_time = $3, duration = $4, recurring_interval = $5, recurring_type = $6`
	for _, schedule := range event.Schedules {
		_, err = tx.ExecContext(ctx, upsertScheduleSql, schedule.ID, event.ID, schedule.StartTime, schedule.DurationInMinutes, schedule.IsFullDay, schedule.RecurringInterval, schedule.RecurringType)
		if err != nil {
			log.Error(err)
		}
	}

	upsertInvitationSql := `INSERT INTO invitation (id, event_id, user_id, token, status) VALUES ($1, $2, $3, $4, $5)
						ON CONFLICT (id, event_id)
						DO UPDATE SET user_id = $3, token = $4, status = $5`
	for _, invitation := range event.Invitations {
		_, err = tx.ExecContext(ctx, upsertInvitationSql, invitation.ID, event.ID, invitation.UserID, invitation.Token, invitation.Status)
		if err != nil {
			log.Error(err)
		}
	}

	return tx.Commit()
}

func (e *EventRepository) FindByID(ctx context.Context, id string) (*core.Event, error) {
	log := logger.WithFields(logger.Fields{
		"id": id,
	})

	var event core.Event
	err := e.dbConn.QueryRowxContext(ctx, `SELECT * FROM event WHERE id = $1`, id).StructScan(&event)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var schedules []core.Schedule
	err = e.dbConn.SelectContext(ctx, &schedules, `SELECT * FROM schedule WHERE event_id = $1`, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	event.Schedules = schedules

	var invitations []core.Invitation
	err = e.dbConn.SelectContext(ctx, &invitations, `SELECT * FROM invitation WHERE event_id = $1`, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	event.Invitations = invitations

	return &event, err
}
