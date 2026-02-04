package models

type ApiResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type PaginatedParams struct {
	Data       any `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type Balance struct {
	Balance  float64 `json:"balance"`
	Currency float64 `json:"currency"`
}
