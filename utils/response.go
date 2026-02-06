package utils

import (
	"encoding/json"
	"net/http"

	"github.com/wizzyszn/go_bank/models"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteSuccess(w http.ResponseWriter, data any) error {
	response := models.ApiResponse{
		Success: true,
		Data:    data,
	}
	return WriteJSON(w, http.StatusOK, response)
}

func WriteCreated(w http.ResponseWriter, data any) error {
	response := models.ApiResponse{
		Success: true,
		Data:    data,
		Message: "Resource created successfully",
	}

	return WriteJSON(w, http.StatusCreated, response)
}

func WriteError(w http.ResponseWriter, status int, message string) error {
	response := models.ApiResponse{
		Error:   message,
		Success: false,
	}
	return WriteJSON(w, status, response)
}

func WriteBadRequest(w http.ResponseWriter, message string) error {

	return WriteError(w, http.StatusBadRequest, message)
}

func WriteUnAuthorized(w http.ResponseWriter, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return WriteError(w, http.StatusUnauthorized, message)
}

func WriteForbidden(w http.ResponseWriter, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return WriteError(w, http.StatusForbidden, message)
}

func WriteNotFound(w http.ResponseWriter, message string) error {
	if message == "" {
		message = "Resource Not Found"
	}
	return WriteError(w, http.StatusNotFound, message)
}

func WriteInternalError(w http.ResponseWriter, message string) error {
	if message == "" {
		message = ""
	}
	return WriteError(w, http.StatusInternalServerError, message)
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return ErrEmptyBody
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

var (
	ErrEmptyBody = &ValidationError{Field: "body", Message: "request body cannot be empty"}
)
