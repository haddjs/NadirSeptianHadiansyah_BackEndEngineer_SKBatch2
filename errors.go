package main

import (
	"encoding/json"
	"net/http"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error      string       `json:"error"`
	StatusCode int          `json:"statusCode"`
	Details    []FieldError `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string, details ...FieldError) {
	writeJSON(w, status, ErrorResponse{Error: msg, StatusCode: status, Details: details})
}
