package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type AuthHandler struct {
	managerService services.ManagerService
	getJWT         utils.GetJWT
	pwdComparator  utils.ComparePassword
}

func NewAuthHandler(
	managerService services.ManagerService,
	getJWT utils.GetJWT,
	pwdComparator utils.ComparePassword,
) *AuthHandler {
	return &AuthHandler{managerService: managerService, getJWT: getJWT, pwdComparator: pwdComparator}
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credential utils.Credential

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cfg, ok := r.Context().Value(constants.ConfigKey).(*config.Config)
	if !ok {
		log.Println("Configuration not found")
		utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	switch credential.Action {
	case utils.Register:
		manager_id, sqlErr := h.managerService.Create(credential.Email, credential.Password)
		if sqlErr != nil {
			switch sqlErr.Type {
			case utils.SQLUniqueViolated:
				utils.SendErrorResponse(w, "Email already registered", http.StatusConflict)
				return
			case utils.InvalidEmailFormat:
				utils.SendErrorResponse(w, "Invalid email format", http.StatusBadRequest)
				return
			case utils.InvalidPasswordLength:
				utils.SendErrorResponse(w, "Invalid password length (min length: 8, max length: 32)", http.StatusBadRequest)
				return
			default:
				log.Println("Failed to create manager:", err)
				utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		token, err := h.getJWT(cfg.JWT.Secret, manager_id, credential.Email)
		if err != nil {
			log.Println("Failed to generate JWT:", err)
			utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := utils.AuthResponse{
			Email: credential.Email,
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to send response to %s: %v", credential.Email, err)
			utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	case utils.Login:
		manager, sqlErr := h.managerService.GetByEmail(credential.Email)
		if sqlErr != nil {
			utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
			return
		}

		error := h.pwdComparator([]byte(manager.Password), []byte(credential.Password))
		if error != nil {
			utils.SendErrorResponse(w, "Invalid credential", http.StatusUnauthorized)
			return
		}

		token, err := h.getJWT(cfg.JWT.Secret, manager.ID, manager.Email)
		if err != nil {
			log.Printf("Failed to generate JWT for user %d: %v", manager.ID, err)
			utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := utils.AuthResponse{
			Email: manager.Email,
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Failed to send response to %d: %v", manager.ID, err)
			utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	default:
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
}
