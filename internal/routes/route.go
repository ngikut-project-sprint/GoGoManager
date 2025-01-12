package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/database"
	"github.com/ngikut-project-sprint/GoGoManager/internal/handlers"
	"github.com/ngikut-project-sprint/GoGoManager/internal/middleware"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repository"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

func NewRouter(cfg *config.Config, db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	ManagerRouter(mux, cfg, db)
	EmployeeRouter(mux, cfg, db)
	return mux
}
func ManagerRouter(mux *http.ServeMux, cfg *config.Config, db *sql.DB) {
	dbAdapter := &database.SqlDBAdapter{DB: db}
	repo := repository.NewManagerRepository(dbAdapter, bcrypt.GenerateFromPassword)
	service := services.NewManagerService(repo, validators.ValidateEmail, validators.ValidatePassword)
	AuthRouter(mux, cfg, service)
	ManagersRouter(mux, cfg, service)
}

func ManagersRouter(mux *http.ServeMux, cfg *config.Config, manager_service services.ManagerService) {
	handler := handlers.NewManagerHandler(manager_service)
	mux.Handle("/v1/user", middleware.ConfigMiddleware(cfg, middleware.AuthMiddleware(jwt.ParseWithClaims, http.HandlerFunc(handler.Manager))))
}

func EmployeeRouter(mux *http.ServeMux, cfg *config.Config, db *sql.DB) {
	repo := repository.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	// Handle /v1/employee for GET (list) and POST (create)
	mux.Handle("/v1/employee", middleware.ConfigMiddleware(cfg,
		middleware.AuthMiddleware(jwt.ParseWithClaims, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handler.List(w, r)
			case http.MethodPost:
				handler.Create(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	))

	// Handle /v1/employee/{identityNumber} for PATCH and DELETE
	mux.Handle("/v1/employee/", middleware.ConfigMiddleware(cfg,
		middleware.AuthMiddleware(jwt.ParseWithClaims, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract identityNumber from path
			identityNumber := strings.TrimPrefix(r.URL.Path, "/v1/employee/")
			if identityNumber == "" {
				http.Error(w, "Missing employee identity number", http.StatusBadRequest)
				return
			}

			switch r.Method {
			case http.MethodPatch:
				handler.Update(w, r, identityNumber)
			case http.MethodDelete:
				handler.Delete(w, r, identityNumber)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})),
	))
}

func AuthRouter(mux *http.ServeMux, cfg *config.Config, manager_service services.ManagerService) {
	handler := handlers.NewAuthHandler(manager_service, utils.GenerateJWT, bcrypt.CompareHashAndPassword)
	mux.Handle("/v1/auth", middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)))
	mux.Handle("/v1/protected", middleware.ConfigMiddleware(cfg, middleware.AuthMiddleware(jwt.ParseWithClaims, http.HandlerFunc(handlers.ExampleSecureHander))))
}
