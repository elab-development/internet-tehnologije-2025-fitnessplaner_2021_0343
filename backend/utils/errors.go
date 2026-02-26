package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse predstavlja JSON odgovor sa greškom
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// JSONError šalje JSON odgovor sa greškom
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}


