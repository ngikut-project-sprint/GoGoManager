package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/handlers"
	"github.com/ngikut-project-sprint/GoGoManager/internal/middleware"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repositories"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
)

func NewRouter(cfg *config.Config, db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	ManagerRouter(mux, cfg, db)
	EmployeeRouter(mux, cfg, db)
	return mux
}

func ManagerRouter(mux *http.ServeMux, cfg *config.Config, db *sql.DB) {
	repo := repositories.NewManagerRepository(db)
	service := services.NewManagerService(repo)
	AuthRouter(mux, cfg, service)
}

func EmployeeRouter(mux *http.ServeMux, cfg *config.Config, db *sql.DB) {
	repo := repositories.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	// Handle /v1/employee for GET (list) and POST (create)
	mux.Handle("/v1/employee", middleware.ConfigMiddleware(cfg,
		middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// Handle /v1/employee/{identityNumber} for PATCH
	mux.Handle("/v1/employee/", middleware.ConfigMiddleware(cfg,
		middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Extract identityNumber from path
			identityNumber := strings.TrimPrefix(r.URL.Path, "/v1/employee/")
			if identityNumber == "" {
				http.Error(w, "Missing employee identity number", http.StatusBadRequest)
				return
			}

			handler.Update(w, r, identityNumber)
		})),
	))
}

func AuthRouter(mux *http.ServeMux, cfg *config.Config, manager_service services.ManagerService) {
	handler := handlers.NewAuthHandler(manager_service)
	mux.Handle("/auth", middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)))
	mux.Handle("/protected", middleware.ConfigMiddleware(cfg, middleware.AuthMiddleware(http.HandlerFunc(handlers.ExampleSecureHander))))
}
