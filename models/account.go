package models

import "time"

// Account represents a Bank Account

type Account struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FirstName    string    `json:"first_name" db:"fisrt_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Balance      float64   `json:"balance" db:"balance"`
	Currency     string    `json:"currency" db:"currency"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateAccountRequest represents the request body for creating an account

type CreateAccountRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

// LoginRequest represents the request body for login

type LoginAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateAccountRequest represents the request body for updating account details

type UpdateAccountRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password,omitempty"`
}

// AccountResponse is what we return to the client (without sensitive data)
type AccountResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Balance   float64   `json:"balance"`
	Status    string    `json:"status"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts Account to AccountResponse (removes sensitive fields)
func (a *Account) ToResponse() *AccountResponse {
	return &AccountResponse{
		ID:        a.ID,
		Email:     a.Email,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Balance:   a.Balance,
		Status:    a.Status,
		Currency:  a.Currency,
		CreatedAt: a.CreatedAt,
	}
}

const (
	AccountStatusActice    string = "active"
	AccountStatusSuspended string = "suspended"
	AccountStatusClosed    string = "close"
)
