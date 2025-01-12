package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes a JSON response with the given status code and data
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        SendErrorResponse(w, "Failed create department", http.StatusBadRequest)
        return
    }
}

// For general response
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
