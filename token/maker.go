package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid or not and in positive case returns token claims
	VerifyToken(token string) (*Payload, error)
}
