package web

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response format.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *ApiError  `json:"error,omitempty"`
}

// ApiError is the standard API error format.
type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RespondWithError sends an error response.
func RespondWithError(w http.ResponseWriter, code string, message string, statusCode int) {
	RespondWithJSON(w, statusCode, Response{
		Error: &ApiError{
			Code:    code,
			Message: message,
		},
	})
}

// RespondWithJSON sends a JSON response.
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(response)
}
