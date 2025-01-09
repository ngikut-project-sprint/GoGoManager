package routes

import (
	"database/sql"
	"fmt"
	"net/http"

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
	// Initialize repository, service, and handler
	repo := repositories.NewEmployeeRepository(db)
	service := services.NewEmployeeService(repo)
	handler := handlers.NewEmployeeHandler(service)

	// Register routes without auth middleware for GET
	mux.HandleFunc("/v1/employee", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Iqbal ganteng")
		if r.Method == http.MethodGet {
			handler.List(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}

func AuthRouter(mux *http.ServeMux, cfg *config.Config, manager_service services.ManagerService) {
	handler := handlers.NewAuthHandler(manager_service)
	mux.Handle("/auth", middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)))
	mux.Handle("/protected", middleware.ConfigMiddleware(cfg, middleware.AuthMiddleware(http.HandlerFunc(handlers.ExampleSecureHander))))
}
