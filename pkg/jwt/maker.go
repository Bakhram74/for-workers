package jwt

import (
	"time"

	db "github.com/ShamilKhal/shgo/db/sqlc"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific user and duration
	CreateToken(user db.User, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
