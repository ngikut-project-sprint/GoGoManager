package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
	"github.com/ngikut-project-sprint/GoGoManager/internal/models"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type ManagerHandler struct {
	managerService services.ManagerService
}

func NewManagerHandler(managerService services.ManagerService) *ManagerHandler {
	return &ManagerHandler{managerService: managerService}
}
func (h *ManagerHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		http.Error(w, "User not aunthenticated", http.StatusUnauthorized)
		return
	}

	manager, err := h.managerService.GetByID(claims.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := manager.ToManagerResponse()

	jsonResponse, error := json.Marshal(response)
	if error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResponse); err != nil {
		log.Println("Failed to write response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *ManagerHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var input models.Manager
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		http.Error(w, "User not aunthenticated", http.StatusUnauthorized)
		return
	}

	manager, err := h.managerService.GetByID(claims.ID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	//assign manager id
	input.ID = manager.ID

	//update manager
	updateErr := h.managerService.Update(&input)
	if updateErr != nil {
		http.Error(w, "Failed to update manager: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//get updated manager
	result, queryErr := h.managerService.GetByID(claims.ID)
	if queryErr != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := result.ToManagerResponse()
	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
func (h *ManagerHandler) Manager(w http.ResponseWriter, r *http.Request) {
	http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.GetUser(w, r) // Handle GET requests
			case http.MethodPatch:
				h.UpdateUser(w, r) // Handle PATCH requests
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}).ServeHTTP(w, r)
}
