package internal

import (
	"time"

	"github.com/satori/uuid"
)

type RecurringType uint

const (
	RecurringType_None RecurringType = iota
	RecurringType_Daily
	RecurringType_Every_Week
	RecurringType_Every_Month
)

type Schedule struct {
	ID            string        `validate:"required"`
	EventID       string        `validate:"required"`
	StartTime     time.Time     `validate:"required"`
	Duration      time.Duration `validate:"required"`
	IsFullDay     bool          `validate:"required"`
	RecurringType RecurringType `validate:"required"`
}

func (s *Schedule) EndTime() time.Time {
	return s.StartTime.Add(s.Duration)
}

func NewSchedule(eventID string, start, end string, isFullDay bool) (Schedule, error) {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return Schedule{}, err
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return Schedule{}, err
	}

	return Schedule{
		ID:        uuid.NewV4().String(),
		EventID:   eventID,
		StartTime: startTime.UTC(),
		Duration:  endTime.UTC().Sub(startTime.UTC()),
		IsFullDay: isFullDay,
	}, nil
}

func (s *Schedule) RecurringInterval() int64 {
	switch s.RecurringType {
	case RecurringType_Daily:
		return int64(time.Duration(24 * time.Hour).Seconds())
	case RecurringType_Every_Week:
		return int64(time.Duration(24 * 7 * time.Hour).Seconds())
	case RecurringType_Every_Month:
		return 0
	default:
		return 0
	}
}
