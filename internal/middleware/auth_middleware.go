package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ngikut-project-sprint/GoGoManager/internal/config"
	"github.com/ngikut-project-sprint/GoGoManager/internal/constants"
	"github.com/ngikut-project-sprint/GoGoManager/internal/utils"
)

func AuthMiddleware(parseJWT utils.ParseJWT, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenRaw := strings.Split(authHeader, " ")
		if len(tokenRaw) != 2 {
			utils.SendErrorResponse(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenString := tokenRaw[1]

		// Get jwt secret from config
		cfg, ok := r.Context().Value(constants.ConfigKey).(*config.Config)
		if !ok {
			log.Println("Configuration not found")
			utils.SendErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Parse token to jwt
		token, err := parseJWT(tokenString, &utils.Claims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			utils.SendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Get jwt claims
		claims, ok := token.Claims.(*utils.Claims)
		if !ok {
			utils.SendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		expirationTime := claims.ExpiresAt.Time
		if time.Now().After(expirationTime) {
			utils.SendErrorResponse(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), constants.JWTKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
