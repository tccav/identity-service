package entities

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID             string
	UserID         string
	ExpirationDate time.Time
	Hash           string
}

func NewToken(userID string, expirationDate time.Time) Token {
	return Token{
		ID:             uuid.NewString(),
		UserID:         userID,
		ExpirationDate: expirationDate,
	}
}
