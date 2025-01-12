package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
    Details string `json:"details,omitempty"`
}

type errorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func SendErrorResponse(w http.ResponseWriter, msg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := errorResponse{
		ErrorMessage: msg,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, msg, statusCode)
		return
	}
}
