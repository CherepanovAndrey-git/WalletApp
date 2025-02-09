package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
	"wallet-app/internal/database"
	"wallet-app/internal/models/apicfg"
	"wallet-app/internal/utils"
)

// @Summary Register a new user
// @Description Create a new user with the provided name
// @Tags users
// @Accept json
// @Produce json
// @Param user body parameters true "User registration details"
// @Success 200 {object} database.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]

type apiConfig struct {
	DB *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

func HandlerCreateUser(apiCfg *apicfg.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Name string `json:"name"`
		}
		var params parameters

		// Decode the request body into the parameters struct
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		// Create the user in the database
		user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
		})
		if err != nil {
			utils.RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
			return
		}

		// Respond with the created user
		utils.RespondWithJSON(w, 200, databaseUserToUser(user))
	}
}

func HandlerGetUser(apiCfg *apicfg.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the user from the context (set by the middleware)
		user := r.Context().Value("user").(database.User)

		// Respond with the user details
		utils.RespondWithJSON(w, 200, databaseUserToUser(user))
	}
}
