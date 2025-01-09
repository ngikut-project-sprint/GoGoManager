package utils

import (
    "encoding/json"
    "net/http"
)

// Response is the standard API response structure
type Response struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Code    string      `json:"code"`
    Details string      `json:"details,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

// WriteJSON writes a JSON response with the given status code and data
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

