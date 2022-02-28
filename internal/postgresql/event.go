package postgresql

import (
	"context"

	"github.com/dzakaammar/event-scheduling-example/internal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (e *EventRepository) Store(ctx context.Context, event *internal.Event) error {
	trx := e.db.Begin()

	err := trx.WithContext(ctx).Create(event).Error
	if err != nil {
		trx.Rollback()
		return err
	}

	return trx.Commit().Error
}

func (e *EventRepository) DeleteByID(ctx context.Context, id string) error {
	err := e.db.WithContext(ctx).Where("id = ?", id).Delete(&internal.Event{}).Error
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (e *EventRepository) Update(ctx context.Context, event *internal.Event) error {
	trx := e.db.Begin()

	err := trx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Updates(event).Error
	if err != nil {
		trx.Rollback()
		return err
	}

	return trx.Commit().Error
}

func (e *EventRepository) FindByID(ctx context.Context, id string) (*internal.Event, error) {
	var event *internal.Event

	err := e.db.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}

	return event, nil
}
