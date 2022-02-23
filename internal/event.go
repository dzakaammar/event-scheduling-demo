package internal

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/satori/uuid"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Event struct {
	ID          string
	Title       string       `validate:"required"`
	Description string       `validate:"required"`
	Timezone    string       `validate:"required"`
	Schedules   []Schedule   `validate:"required,dive,required"`
	Invitations []Invitation `validate:"dive"`
	CreatedBy   string       `validate:"required"`
	CreatedAt   time.Time    `validate:"required"`
	UpdatedAt   *time.Time
}

func (e *Event) Validate() error {
	return validate.Struct(e)
}

func NewEvent() *Event {
	return &Event{
		ID:        uuid.NewV4().String(),
		CreatedAt: time.Now(),
	}
}

//go:generate mockgen -destination=mock/mock_event_repository.go -package=mock github.com/dzakaammar/event-scheduling-example/internal EventRepository
type EventRepository interface {
	Store(ctx context.Context, e *Event) error
	DeleteByID(ctx context.Context, id string) error
	Update(ctx context.Context, e *Event) error
	FindByID(ctx context.Context, id string) (*Event, error)
}
