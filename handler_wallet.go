package main

import (
    "database/sql"
    "encoding/json"
    _"fmt"
    "log"
    "net/http"
	"time"

    "github.com/google/uuid"
)

type Wallet struct {
    ID        int       `json:"id"`
    UUID      string    `json:"uuid"`
    UserID    string       `json:"user_id"`
    Balance   float64   `json:"balance"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CreateWalletRequest struct {
    UserID int `json:"user_id"`
}
// @Summary Create a new wallet
// @Description Create a new wallet for the authenticated user
// @Tags wallets
// @Accept json
// @Produce json
// @Param ApiKey header string true "API Key"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallets [post]
func CreateWalletHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract API key from the ApiKey header
        apiKey := r.Header.Get("ApiKey")
        if apiKey == "" {
            http.Error(w, "ApiKey header is required", http.StatusUnauthorized)
            return
        }

        // Validate the API key
        userID, err := validateAPIKey(db, apiKey)
        if err != nil {
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
            return
        }

        // Generate a new wallet UUID
        walletUUID := uuid.New().String()

        // Insert the wallet into the database
        query := `INSERT INTO wallets (uuid, user_id, balance) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
        var wallet Wallet
        wallet.UUID = walletUUID
        wallet.UserID = userID // Now userID is a string
        wallet.Balance = 0.0

        err = db.QueryRow(query, wallet.UUID, wallet.UserID, wallet.Balance).Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)
        if err != nil {
            http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
            return
        }

        // Return the created wallet
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{
            "walletId": wallet.UUID,
        })
    }
}

// validateAPIKey checks if the API key is valid and returns the associated user ID
func validateAPIKey(db *sql.DB, apiKey string) (string, error) {
    var userID string
    err := db.QueryRow("SELECT id FROM users WHERE api_key = $1", apiKey).Scan(&userID)
    if err != nil {
        log.Printf("Error validating API key: %v", err)
        log.Printf("Provided API key: %s", apiKey)
        return "", err
    }
    log.Printf("API key validated for user ID: %s", userID)
    return userID, nil
}