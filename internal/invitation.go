package internal

import (
	"encoding/base64"

	"github.com/satori/uuid"
)

type InvitationStatus uint

const (
	InvitationStatus_Unknown InvitationStatus = iota
	InvitationStatus_Confirmed
	InvitationStatus_Declined
)

type Invitation struct {
	ID      string           `validate:"required"`
	EventID string           `validate:"required"`
	UserID  string           `validate:"required"`
	Status  InvitationStatus `validate:"required"`
	Token   string           `validate:"required"`
}

func NewInvitation(eventID string, userID string) *Invitation {
	id := uuid.NewV4().String()
	return &Invitation{
		ID:      id,
		EventID: eventID,
		UserID:  userID,
		Status:  InvitationStatus_Unknown,
		Token:   base64.StdEncoding.EncodeToString([]byte(id)), // use base64-encoded id for simplicity sake
	}
}
