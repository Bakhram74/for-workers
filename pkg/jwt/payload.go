package jwt

import (
	"errors"

	"time"

	db "github.com/ShamilKhal/shgo/db/sqlc"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role,omitempty"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific user and duration
func NewPayload(user db.User, duration time.Duration) *Payload {

	payload := &Payload{
		ID:        user.ID,
		Name:      user.Name,
		Phone:     user.Phone,
		Role:      user.Role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error { //TODO delete?
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
