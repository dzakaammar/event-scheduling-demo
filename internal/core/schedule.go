package core

import (
	"fmt"
	"time"

	"github.com/satori/uuid"
)

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
	ID                string        `validate:"required" gorm:"primaryKey"`
	EventID           string        `validate:"required"`
	StartTime         int64         `validate:"required"`
	Duration          time.Duration `validate:"required"`
	IsFullDay         bool
	RecurringType     RecurringType
	RecurringInterval int64
}

func (s *Schedule) TableName() string {
	return "schedule"
}

func (s *Schedule) StartTimeIn(loc string) (time.Time, error) {
	l, err := time.LoadLocation(loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %s", loc)
	}
	return time.Unix(s.StartTime, 0).In(l), nil
}

func (s *Schedule) EndTimeFrom(st time.Time) time.Time {
	return st.Add(s.Duration * time.Minute)
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
		ID:            uuid.NewV4().String(),
		EventID:       eventID,
		StartTime:     startTime.UTC().Unix(),
		Duration:      time.Duration(endTime.UTC().Sub(startTime.UTC()).Minutes()),
		IsFullDay:     isFullDay,
		RecurringType: rt,
	}
	s.RecurringInterval = rt.interval()

	return s, nil
}
