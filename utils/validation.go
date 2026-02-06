package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidateEmail checks if an email is valid
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}

	// Use Go's built-in email parser
	_, err := mail.ParseAddress(email)
	if err != nil {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}

	return nil
}

// ValidateName checks if a name is valid
func ValidateName(name, fieldName string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s is required", fieldName)}
	}

	if len(name) < 2 {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be at least 2 characters", fieldName)}
	}

	if len(name) > 100 {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be at most 100 characters", fieldName)}
	}

	// Only allow letters, spaces, hyphens, and apostrophes
	validName := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	if !validName.MatchString(name) {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s contains invalid characters", fieldName)}
	}

	return nil
}

// ValidateAmount checks if a monetary amount is valid
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return &ValidationError{Field: "amount", Message: "amount must be greater than 0"}
	}

	// Check for reasonable maximum (e.g., $1 million per transaction)
	if amount > 1000000 {
		return &ValidationError{Field: "amount", Message: "amount exceeds maximum allowed"}
	}

	// Check for too many decimal places (max 2 for cents)
	// This is a simple check - in production you'd use decimal type
	if !isValidMoneyFormat(amount) {
		return &ValidationError{Field: "amount", Message: "amount can have at most 2 decimal places"}
	}

	return nil
}

// isValidMoneyFormat checks if amount has at most 2 decimal places
func isValidMoneyFormat(amount float64) bool {
	// Multiply by 100 and check if it's a whole number
	cents := amount * 100
	return cents == float64(int(cents))
}

// ValidateAccountID checks if an account ID is valid
func ValidateAccountID(id int) error {
	if id <= 0 {
		return &ValidationError{Field: "account_id", Message: "invalid account ID"}
	}
	return nil
}

// ValidatePositiveInt checks if an integer is positive
func ValidatePositiveInt(value int, fieldName string) error {
	if value <= 0 {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s must be positive", fieldName)}
	}
	return nil
}

// ValidatePagination validates pagination parameters
func ValidatePagination(page, limit int) error {
	if page < 1 {
		return &ValidationError{Field: "page", Message: "page must be at least 1"}
	}

	if limit < 1 {
		return &ValidationError{Field: "limit", Message: "limit must be at least 1"}
	}

	if limit > 100 {
		return &ValidationError{Field: "limit", Message: "limit cannot exceed 100"}
	}

	return nil
}

// ValidateRequired checks if a string field is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{Field: fieldName, Message: fmt.Sprintf("%s is required", fieldName)}
	}
	return nil
}

// SanitizeString trims whitespace and removes extra spaces
func SanitizeString(s string) string {
	// Trim leading and trailing whitespace
	s = strings.TrimSpace(s)

	// Replace multiple spaces with single space
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")

	return s
}
