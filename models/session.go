package models

import "time"

type Session struct {
	ID        string    `json:"id" db:"id"`
	AccountID int       `json:"account_id" db:"account_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LoginResponse struct {
	Account   *AccountResponse `json:"account"`
	SessionID string           `json:"session_id"`
	ExpiresAt time.Time        `json:"expires_at"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
