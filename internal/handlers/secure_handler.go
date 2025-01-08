package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

func ExampleSecureHander(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		http.Error(w, "User not aunthenticated", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"id":    claims.ID,
		"email": claims.Email,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Failed to marshal response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(jsonResponse); err != nil {
		log.Println("Failed to write response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
