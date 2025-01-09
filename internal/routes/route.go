package routes

import (
	"database/sql"
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
	DepartmentRouter(mux, cfg, db)
	return mux
}

func ManagerRouter(mux *http.ServeMux, cfg *config.Config, db *sql.DB) {
	repo := repositories.NewManagerRepository(db)
	service := services.NewManagerService(repo)
	AuthRouter(mux, cfg, service)
}

func AuthRouter(mux *http.ServeMux, cfg *config.Config, manager_service services.ManagerService) {
	handler := handlers.NewAuthHandler(manager_service)
	mux.Handle("/auth", middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)))
	mux.Handle("/protected", middleware.ConfigMiddleware(cfg, middleware.AuthMiddleware(http.HandlerFunc(handlers.ExampleSecureHander))))
}

func DepartmentRouter(mux *http.ServeMux, cfg *config.Config, db *sql.DB) {
    repo := repositories.NewDepartmentRepository(db)
    service := services.NewDepartmentService(repo)
    handler := handlers.NewDepartmentHandler(service)

    mux.Handle("/department", middleware.ConfigMiddleware(cfg, 
        middleware.AuthMiddleware(http.HandlerFunc(handler.HandleDepartment))))
    mux.Handle("/department/", middleware.ConfigMiddleware(cfg, 
        middleware.AuthMiddleware(http.HandlerFunc(handler.HandleDepartmentWithID))))
}