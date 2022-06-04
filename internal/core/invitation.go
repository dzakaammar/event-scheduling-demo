package core

import (
	"encoding/base64"
	"time"

	"github.com/satori/uuid"
)

type InvitationStatus uint

const (
	InvitationStatus_Unknown InvitationStatus = iota
	InvitationStatus_Confirmed
	InvitationStatus_Declined
)

type Invitation struct {
	ID        string           `validate:"required" db:"id"`
	EventID   string           `validate:"required" db:"event_id"`
	UserID    string           `validate:"required" db:"user_id"`
	Status    InvitationStatus `db:"status"`
	Token     string           `validate:"required" db:"token"`
	UpdatedAt *time.Time       `db:"updated_at"`
}

func NewInvitation(eventID string, userID string) Invitation {
	id := uuid.NewV4().String()
	return Invitation{
		ID:      id,
		EventID: eventID,
		UserID:  userID,
		Status:  InvitationStatus_Unknown,
		Token:   base64.StdEncoding.EncodeToString([]byte(id)), // use base64-encoded id for simplicity sake
	}
}
