package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"wallet-app/internal/database"
	"wallet-app/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletOperationRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// @Summary Create a new wallet
// @Description Create a new wallet for authenticated user
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]interface{} "Wallet created successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/create-wallet [post]

func CreateWalletHandler(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		// Check if wallet already exists
		existingWallet, err := db.GetWalletByUserID(r.Context(), userID)
		if err == nil {
			
			utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"message": "Wallet already exists",
				"balance": existingWallet.Balance,
			})
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check wallet")
			return
		}

		// Create new wallet only if it doesn't exist
		wallet, err := db.CreateWallet(r.Context(), userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create wallet")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Wallet created successfully",
			"balance": wallet.Balance,
		})
	}
}

// @Summary Deposit money
// @Description Deposit money into user's wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body WalletOperationRequest true "Deposit details"
// @Success 200 {object} map[string]interface{} "Operation successful"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/wallet/deposit [post]
func WalletOperationHandler(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		var req WalletOperationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		
		isDeposit := strings.HasSuffix(r.URL.Path, "/deposit")
		amount := req.Amount
		if !isDeposit {
			amount = -amount 
		}

		
		err := db.UpdateWalletBalance(r.Context(), database.UpdateWalletBalanceParams{
			Userid: userID,
			Amount: fmt.Sprintf("%.2f", amount),
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update balance")
			return
		}

		
		updatedWallet, err := db.GetWalletByUserID(r.Context(), userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch updated balance")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Operation successful",
			"balance": updatedWallet.Balance,
		})
	}
}

// @Summary Get wallet balance
// @Description Get current balance of user's wallet
// @Tags wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]float64 "Wallet balance"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Wallet not found"
// @Router /v1/balance [get]
func GetBalanceHandler(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		wallet, err := db.GetWalletByUserID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.RespondWithError(w, http.StatusNotFound, "Wallet not found")
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch wallet")
			}
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]float64{
			"balance": utils.ParseStringToFloat64(wallet.Balance),
		})
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
