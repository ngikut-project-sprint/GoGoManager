package utils

import (
    "net/http"
)

// Common error responses
func NotFound(w http.ResponseWriter, message string) {
    SendErrorResponse(w, message, http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, message string) {
    SendErrorResponse(w, message, http.StatusBadRequest)
}

func MethodNotAllowed(w http.ResponseWriter, method string) {
    SendErrorResponse(w,  "Method " + method + " not allowed", http.StatusMethodNotAllowed)
}

func InternalServerError(w http.ResponseWriter, err error) {
    SendErrorResponse(w,  "Internal server error", http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter) {
    SendErrorResponse(w,  "Unauthorized access", http.StatusUnauthorized)
}
