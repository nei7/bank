package token

import "time"

// Interface for manage tokens
type Maker interface {
	// Create token for specific username
	CreateToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
