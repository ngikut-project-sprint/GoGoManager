package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
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
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Println("Failed to write response:", err)
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *ManagerHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var input utils.ManagerRequest
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(constants.JWTKey).(*utils.Claims)
	if !ok {
		http.Error(w, "User not aunthenticated", http.StatusUnauthorized)
		return
	}

	//assign manager id
	input.ID = claims.ID

	//update manager
	updateErr := h.managerService.Update(&input)
	if updateErr != nil {
		switch updateErr.Type {
		case utils.SQLUniqueViolated:
			utils.SendErrorResponse(w, "Email already registered", http.StatusConflict)
			return
		case utils.InvalidEmailFormat:
			utils.SendErrorResponse(w, "Invalid email format", http.StatusBadRequest)
			return
		case utils.InvalidPasswordLength:
			utils.SendErrorResponse(w, "Invalid password length (min length: 8, max length: 32)", http.StatusBadRequest)
			return
		case utils.InvalidNameLength:
			utils.SendErrorResponse(w, "Invalid name length (min length: 4, max length: 52)", http.StatusBadRequest)
			return
		default:
			log.Println("Failed to create manager:", updateErr)
			utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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
	if err := json.NewEncoder(w).Encode(jsonResponse); err != nil {
		log.Println("Failed to write response:", err)
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
