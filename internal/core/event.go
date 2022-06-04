package core

import (
	"context"
	"time"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/go-playground/validator/v10"
	"github.com/satori/uuid"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Event struct {
	ID          string     `validate:"required" db:"id"`
	Title       string     `validate:"required" db:"title"`
	Description string     `validate:"required" db:"description"`
	Timezone    string     `validate:"required" db:"timezone"`
	CreatedBy   string     `db:"created_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`

	Schedules   []Schedule   `validate:"required,dive,required"`
	Invitations []Invitation `validate:"dive"`
}

func (e *Event) Validate() error {
	_, err := time.LoadLocation(e.Timezone)
	if err != nil {
		return internal.WrapErr(internal.ErrInvalidTimezone, e.Timezone)
	}
	return validate.Struct(e)
}

func (e *Event) GetUpdatedAt() string {
	if e.UpdatedAt == nil {
		return time.Time{}.Format(time.RFC3339)
	}

	return e.UpdatedAt.Format(time.RFC3339)
}

func NewEvent(createdBy string) *Event {
	return &Event{
		ID:        uuid.NewV4().String(),
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
	}
}

//go:generate mockgen -destination=../mock/mock_event_repository.go -package=mock github.com/dzakaammar/event-scheduling-example/internal/core EventRepository
type EventRepository interface {
	Store(ctx context.Context, e *Event) error
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, e *Event) error
	FindByID(ctx context.Context, id string) (*Event, error)
}
