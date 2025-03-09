package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal/core"
	"github.com/dzakaammar/event-scheduling-example/internal/postgresql/gen"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	dbConn  *sqlx.DB
	queries *gen.Queries
}

func NewEventRepository(dbConn *sqlx.DB) *EventRepository {
	return &EventRepository{
		dbConn:  dbConn,
		queries: gen.New(dbConn),
	}
}

func (e *EventRepository) Store(ctx context.Context, event *core.Event) error { //nolint:funlen,gocognit
	tx, err := e.dbConn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			slog.Error(rollbackErr.Error())
		}
	}()

	err = e.queries.WithTx(tx).CreateEvent(ctx, gen.CreateEventParams{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Timezone:    event.Timezone,
		CreatedBy:   event.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   sql.NullTime{Time: time.Now()},
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, schedule := range event.Schedules {
		err = e.queries.WithTx(tx).CreateSchedule(ctx, gen.CreateScheduleParams{
			ID:                schedule.ID,
			EventID:           event.ID,
			StartTime:         schedule.StartTime,
			Duration:          schedule.DurationInMinutes,
			IsFullDay:         schedule.IsFullDay,
			RecurringInterval: schedule.RecurringInterval,
			RecurringType:     string(schedule.RecurringType),
		})
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	for _, invitation := range event.Invitations {
		err = e.queries.WithTx(tx).CreateInvitation(ctx, gen.CreateInvitationParams{
			ID:      invitation.ID,
			EventID: invitation.EventID,
			UserID:  int32(invitation.UserID),
			Token:   invitation.Token,
			Status:  int16(invitation.Status), //nolint:gosec
		})
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	return tx.Commit()
}

func (e *EventRepository) DeleteByID(ctx context.Context, id string) error {
	return e.queries.DeleteEvent(ctx, id)
}

func (e *EventRepository) Update(ctx context.Context, event *core.Event) error { //nolint:gocognit
	tx, err := e.dbConn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			slog.Error(rollbackErr.Error())
		}
	}()

	err = e.queries.WithTx(tx).UpdateEvent(ctx, gen.UpdateEventParams{
		Title:       event.Title,
		Description: event.Description,
		Timezone:    event.Timezone,
	})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, schedule := range event.Schedules {
		err = e.queries.WithTx(tx).UpsertSchedule(ctx, gen.UpsertScheduleParams{
			ID:                schedule.ID,
			EventID:           event.ID,
			StartTime:         schedule.StartTime,
			Duration:          schedule.DurationInMinutes,
			IsFullDay:         schedule.IsFullDay,
			RecurringInterval: schedule.RecurringInterval,
			RecurringType:     string(schedule.RecurringType),
		})
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	for _, invitation := range event.Invitations {
		err = e.queries.WithTx(tx).UpsertInvitation(ctx, gen.UpsertInvitationParams{
			ID:      invitation.ID,
			EventID: event.ID,
			UserID:  invitation.UserID,
			Token:   invitation.Token,
			Status:  int16(invitation.Status), //nolint:gosec
		})
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	return tx.Commit()
}

func (e *EventRepository) FindByID(ctx context.Context, id string) (*core.Event, error) {
	queryEvent, err := e.queries.FindEventByID(ctx, id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	event := core.Event{
		ID:          queryEvent.ID,
		Title:       queryEvent.Title,
		Description: queryEvent.Description,
		Timezone:    queryEvent.Timezone,
		CreatedBy:   queryEvent.CreatedBy,
		CreatedAt:   queryEvent.CreatedAt,
		UpdatedAt:   &queryEvent.UpdatedAt.Time,
	}

	var schedules []core.Schedule
	err = e.dbConn.SelectContext(ctx, &schedules, `SELECT * FROM schedule WHERE event_id = $1`, id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	event.Schedules = schedules

	var invitations []core.Invitation
	err = e.dbConn.SelectContext(ctx, &invitations, `SELECT * FROM invitation WHERE event_id = $1`, id)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	event.Invitations = invitations

	return &event, err
}
