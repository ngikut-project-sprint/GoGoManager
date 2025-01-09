package utils

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func SendErrorResponse(w http.ResponseWriter, msg string, statusCode int) {
	response := errorResponse{
		ErrorMessage: msg,
	}
	json.NewEncoder(w).Encode(response)
}
