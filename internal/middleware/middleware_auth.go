package middleware

import (
	"fmt"
	"net/http"
	"wallet-app/internal/auth"
	"wallet-app/internal/database"
	"wallet-app/internal/models/apicfg"
	"wallet-app/internal/utils"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func Auth(apiCfg *apicfg.ApiConfig, handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIkey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't get use: %v", err))
			return
		}

		handler(w, r, user)
	}
}
