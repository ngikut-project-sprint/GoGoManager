package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

type AuthHandler struct {
	managerService services.ManagerService
}

func NewAuthHandler(managerService services.ManagerService) *AuthHandler {
	return &AuthHandler{managerService: managerService}
}

func (h *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var credential utils.Credential

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cfg, ok := r.Context().Value(constants.ConfigKey).(*config.Config)
	if !ok {
		log.Println("Configuration not found")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	switch credential.Action {
	case utils.Register:
		manager_id, sqlErr := h.managerService.Create(credential.Email, credential.Password)
		if sqlErr != nil {
			switch sqlErr.Type {
			case utils.SQLUniqueViolated:
				http.Error(w, "Email already registered", http.StatusConflict)
				return
			case utils.InvalidEmailFormat:
				http.Error(w, "Invalid email format", http.StatusBadRequest)
				return
			case utils.InvalidPasswordLength:
				http.Error(w, "Invalid password length (min length: 8, max length: 32)", http.StatusBadRequest)
				return
			default:
				log.Println("Failed to create manager:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		token, err := utils.GenerateJWT(cfg.JWT.Secret, manager_id, credential.Email)
		if err != nil {
			log.Println("Failed to generate JWT:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := utils.AuthResponse{
			Email: credential.Email,
			Token: token,
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

	case utils.Login:
		manager, sqlErr := h.managerService.GetByEmail(credential.Email)
		if sqlErr != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		error := bcrypt.CompareHashAndPassword([]byte(*manager.Password), []byte(credential.Password))
		if error != nil {
			http.Error(w, "Invalid credential", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(cfg.JWT.Secret, manager.ID, manager.Email)
		if err != nil {
			log.Printf("Failed to generate JWT for user %d: %v", manager.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := utils.AuthResponse{
			Email: manager.Email,
			Token: token,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to marshal response for user %d: %v", manager.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(jsonResponse); err != nil {
			log.Printf("Failed to write response for user %d: %v", manager.ID, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
}
