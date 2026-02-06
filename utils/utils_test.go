package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wizzyszn/go_bank/models"

	
)

// Test Password Utilities
func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mySecurePassword123"
	hash, _ := HashPassword(password)

	// Correct password
	if err := CheckPassword(password, hash); err != nil {
		t.Errorf("Valid password should pass: %v", err)
	}

	// Wrong password
	if err := CheckPassword("wrongPassword", hash); err == nil {
		t.Error("Invalid password should fail")
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		password  string
		shouldErr bool
	}{
		{"short", true},                  // Too short
		{"validPassword123", false},      // Valid
		{"12345678", false},              // Valid (8 chars)
		{"a", true},                      // Too short
		{string(make([]byte, 73)), true}, // Too long
	}

	for _, tt := range tests {
		err := ValidatePasswordStrength(tt.password)
		if tt.shouldErr && err == nil {
			t.Errorf("Password '%s' should fail validation", tt.password)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Password '%s' should pass validation: %v", tt.password, err)
		}
	}
}

// Test Validation Utilities
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email     string
		shouldErr bool
	}{
		{"test@example.com", false},
		{"invalid-email", true},
		{"", true},
		{"test@", true},
		{"@example.com", true},
		{"test.user+tag@example.co.uk", false},
	}

	for _, tt := range tests {
		err := ValidateEmail(tt.email)
		if tt.shouldErr && err == nil {
			t.Errorf("Email '%s' should fail validation", tt.email)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Email '%s' should pass validation: %v", tt.email, err)
		}
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name      string
		shouldErr bool
	}{
		{"John", false},
		{"Mary-Jane", false},
		{"O'Brien", false},
		{"A", true},                       // Too short
		{"", true},                        // Empty
		{"John123", true},                 // Contains numbers
		{"John@Smith", true},              // Invalid characters
		{string(make([]byte, 101)), true}, // Too long
	}

	for _, tt := range tests {
		err := ValidateName(tt.name, "name")
		if tt.shouldErr && err == nil {
			t.Errorf("Name '%s' should fail validation", tt.name)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Name '%s' should pass validation: %v", tt.name, err)
		}
	}
}

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		amount    float64
		shouldErr bool
	}{
		{100.00, false},
		{0.01, false},
		{999999.99, false},
		{0.00, true},       // Zero
		{-50.00, true},     // Negative
		{1000001.00, true}, // Too large
		{100.123, true},    // Too many decimals
	}

	for _, tt := range tests {
		err := ValidateAmount(tt.amount)
		if tt.shouldErr && err == nil {
			t.Errorf("Amount %.2f should fail validation", tt.amount)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Amount %.2f should pass validation: %v", tt.amount, err)
		}
	}
}

// Test Response Utilities
func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"message": "test"}

	err := WriteJSON(w, http.StatusOK, data)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestWriteSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	WriteSuccess(w, data)

	var response models.ApiResponse
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Error != "" {
		t.Error("Expected error to be empty")
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	WriteError(w, http.StatusBadRequest, "test error")

	var response models.ApiResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Error != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", response.Error)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestParseJSON(t *testing.T) {
	// Valid JSON
	body := `{"name":"John","age":30}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(body))

	var data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err := ParseJSON(req, &data)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if data.Name != "John" {
		t.Errorf("Expected name 'John', got '%s'", data.Name)
	}

	// Invalid JSON
	req = httptest.NewRequest("POST", "/test", bytes.NewBufferString("invalid json"))
	err = ParseJSON(req, &data)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Empty body
	req = httptest.NewRequest("POST", "/test", nil)
	err = ParseJSON(req, &data)
	if err == nil {
		t.Error("Expected error for nil body")
	}
}

// Test Session Utilities
func TestGenerateSessionID(t *testing.T) {
	id1, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("Failed to generate session ID: %v", err)
	}

	if id1 == "" {
		t.Error("Session ID should not be empty")
	}

	// Generate another and ensure they're different
	id2, _ := GenerateSessionID()
	if id1 == id2 {
		t.Error("Session IDs should be unique")
	}

	// Check length (32 bytes base64 encoded should be ~43 characters)
	if len(id1) < 40 {
		t.Error("Session ID seems too short")
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"hello   world", "hello world"},
		{"\thello\n", "hello"},
		{"  multiple   spaces  ", "multiple spaces"},
	}

	for _, tt := range tests {
		result := SanitizeString(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeString('%s') = '%s', expected '%s'", tt.input, result, tt.expected)
		}
	}
}
