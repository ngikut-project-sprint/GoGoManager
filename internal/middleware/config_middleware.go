package middleware

import (
	"context"
	"net/http"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
)

func ConfigMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), constants.ConfigKey, cfg)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
