package routes

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/database"
	"github.com/ngikut-project-sprint/GoGoManager/internal/handlers"
	"github.com/ngikut-project-sprint/GoGoManager/internal/middleware"
	"github.com/ngikut-project-sprint/GoGoManager/internal/repository"
	"github.com/ngikut-project-sprint/GoGoManager/internal/services"
	"github.com/ngikut-project-sprint/GoGoManager/internal/validators"
)

func NewRouter(cfg *config.Config, db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	ManagerRouter(mux, cfg, db)
	return mux
}

func ManagerRouter(mux *http.ServeMux, cfg *config.Config, db *sql.DB) {
	dbAdapter := &database.SqlDBAdapter{DB: db}
	repo := repository.NewManagerRepository(dbAdapter, bcrypt.GenerateFromPassword)
	service := services.NewManagerService(repo, validators.ValidateEmail, validators.ValidatePassword)
	AuthRouter(mux, cfg, service)
}

func AuthRouter(mux *http.ServeMux, cfg *config.Config, manager_service services.ManagerService) {
	handler := handlers.NewAuthHandler(manager_service)
	mux.Handle("/auth", middleware.ConfigMiddleware(cfg, http.HandlerFunc(handler.Auth)))
	mux.Handle("/protected", middleware.ConfigMiddleware(cfg, middleware.AuthMiddleware(http.HandlerFunc(handlers.ExampleSecureHander))))
}
