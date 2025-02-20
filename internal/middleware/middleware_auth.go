package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"wallet-app/internal/auth"
	"wallet-app/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		log.Printf("Auth header: %s", authHeader[:10])

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			log.Printf("Invalid auth format: %v", bearerToken)
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		tokenString := bearerToken[1]

		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			log.Printf("Token parsing error: %v", err)
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		log.Printf("Valid token for user: %s", claims.UserID)

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
