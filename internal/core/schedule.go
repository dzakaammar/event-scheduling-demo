package core

import (
	"errors"
	"time"

	"github.com/satori/uuid"
)

var ErrInvalidTimezone = errors.New("invalid timezone")

type RecurringType string

const (
	RecurringType_None       RecurringType = "NONE"
	RecurringType_Daily      RecurringType = "DAILY"
	RecurringType_Every_Week RecurringType = "WEEK"
)

func (r RecurringType) interval() int64 {
	switch r {
	case RecurringType_Daily:
		return int64(24 * time.Hour.Seconds())
	case RecurringType_Every_Week:
		return int64(24 * 7 * time.Hour.Seconds())
	default:
		return 0
	}
}

type Schedule struct {
	ID                string        `db:"id" validate:"required"`
	EventID           string        `db:"event_id" validate:"required"`
	StartTime         int64         `db:"start_time" validate:"required"`
	DurationInMinutes int64         `db:"duration" validate:"required"`
	IsFullDay         bool          `db:"is_full_day"`
	RecurringType     RecurringType `db:"recurring_type"`
	RecurringInterval int64         `db:"recurring_interval"`
}

func (s *Schedule) StartTimeIn(loc string) (time.Time, error) {
	return time.Unix(s.StartTime, 0).In(time.UTC), nil
}

func (s *Schedule) EndTimeFrom(st time.Time) time.Time {
	return st.Add(time.Duration(s.DurationInMinutes) * time.Minute)
}

func NewSchedule(eventID string, start, end string, isFullDay bool, rt RecurringType) (Schedule, error) {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return Schedule{}, err
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return Schedule{}, err
	}

	s := Schedule{
		ID:                uuid.NewV4().String(),
		EventID:           eventID,
		StartTime:         startTime.UTC().Unix(),
		DurationInMinutes: int64(endTime.UTC().Sub(startTime.UTC()).Minutes()),
		IsFullDay:         isFullDay,
		RecurringType:     rt,
	}
	s.RecurringInterval = rt.interval()

	return s, nil
}
