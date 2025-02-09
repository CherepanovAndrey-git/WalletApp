package handlers

import (
	"net/http"
	"wallet-app/internal/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 400, "something went wrong")
}
