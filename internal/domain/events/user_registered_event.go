package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserRegisteredEvent represents the event published when a user registers
type UserRegisteredEvent struct {
	MessageID string    `json:"messageId"`
	IDCitizen int       `json:"idCitizen"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

// NewUserRegisteredEvent creates a new UserRegisteredEvent with a unique message ID
func NewUserRegisteredEvent(idCitizen int, name, email string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		MessageID: uuid.New().String(),
		IDCitizen: idCitizen,
		Name:      name,
		Email:     email,
		Timestamp: time.Now(),
	}
}

// ToJSON converts the event to JSON bytes
func (e *UserRegisteredEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
