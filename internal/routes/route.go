package routes

import (
	"database/sql"
	"net/http"

	"github.com/ngikut-project-sprint/GoGoManager/internal/handlers"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repositories"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	ManagerRouter(mux, db)
	return mux
}

func ManagerRouter(mux *http.ServeMux, db *sql.DB) {
	repo := repositories.NewManagerRepository(db)
	service := services.NewManagerService(repo)
	AuthRouter(mux, service)
}

func AuthRouter(mux *http.ServeMux, manager_service services.ManagerService) {
	handler := handlers.NewAuthHandler(manager_service)
	mux.HandleFunc("/auth", handler.Auth)
}
