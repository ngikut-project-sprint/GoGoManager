package utils

import (
    "encoding/json"
    "net/http"
)

func WriteJSONError(w http.ResponseWriter, status int, message string, code string, details string) {
    errorResponse := ErrorResponse{
        Status:  status,
        Message: message,
        Code:    code,
        Details: details,
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(errorResponse)
}

// Common error responses
func NotFound(w http.ResponseWriter, message string) {
    WriteJSONError(w, http.StatusNotFound, message, "NOT_FOUND", "")
}

func BadRequest(w http.ResponseWriter, message string) {
    WriteJSONError(w, http.StatusBadRequest, message, "BAD_REQUEST", "")
}

func MethodNotAllowed(w http.ResponseWriter, method string) {
    WriteJSONError(w, http.StatusMethodNotAllowed, 
        "Method " + method + " not allowed", 
        "METHOD_NOT_ALLOWED", 
        "Supported methods: GET, POST, PATCH, DELETE")
}

func InternalServerError(w http.ResponseWriter, err error) {
    WriteJSONError(w, http.StatusInternalServerError, 
        "Internal server error", 
        "INTERNAL_ERROR",
        err.Error())
}

func Unauthorized(w http.ResponseWriter) {
    WriteJSONError(w, http.StatusUnauthorized, 
        "Unauthorized access", 
        "UNAUTHORIZED", 
        "Missing or invalid authentication token")
}
