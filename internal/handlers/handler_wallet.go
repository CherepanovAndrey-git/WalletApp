package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"wallet-app/internal/database"
	"wallet-app/internal/utils"
)

// WalletResponseSwag represents wallet operation response
type WalletResponseSwag struct {
	Message    string             `json:"message" example:"Deposit successful"`
	NewBalance map[string]float64 `json:"new_balance" example:"USD:1000.50,RUB:5000.00,EUR:300.00"`
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
// @Success 201 {object} handlers.WalletResponseSwag "Wallet created successfully"
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

		existingWallet, err := db.GetWalletByUserID(r.Context(), userID)
		if err == nil {
			utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"message": "Wallet already exists",
				"balances": map[string]float64{
					"USD": utils.ParseStringToFloat64(existingWallet.BalanceUsd),
					"RUB": utils.ParseStringToFloat64(existingWallet.BalanceRub),
					"EUR": utils.ParseStringToFloat64(existingWallet.BalanceEur),
				},
			})
			return
		} else if !errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check wallet")
			return
		}

		wallet, err := db.CreateWallet(r.Context(), userID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create wallet")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "Wallet created successfully",
			"balances": map[string]float64{
				"USD": utils.ParseStringToFloat64(wallet.BalanceUsd),
				"RUB": utils.ParseStringToFloat64(wallet.BalanceRub),
				"EUR": utils.ParseStringToFloat64(wallet.BalanceEur),
			},
		})
	}
}

// @Summary Deposit/Withdraw money
// @Description Handle deposit or withdrawal for user's wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param operation path string true "Operation type" Enums(deposit, withdraw)
// @Param request body models.WalletOperationRequestSwag true "Operation details"
// @Success 200 {object} handlers.WalletResponseSwag "Operation successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/wallet/{operation} [post]

func WalletOperationHandler(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		// Get the operation type from the URL path
		operation := chi.URLParam(r, "operation")
		if operation != "deposit" && operation != "withdraw" {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid operation. Must be 'deposit' or 'withdraw'")
			return
		}

		var req WalletOperationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		validCurrencies := map[string]bool{"USD": true, "RUB": true, "EUR": true}
		if !validCurrencies[req.Currency] {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid currency. Allowed: USD, RUB, EUR")
			return
		}

		if req.Amount <= 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "Amount must be positive")
			return
		}

		isDeposit := operation == "deposit"
		operationAmount := req.Amount
		if !isDeposit {
			operationAmount = -operationAmount
		}

		if !isDeposit {
			wallet, err := db.GetWalletByUserID(r.Context(), userID)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to check balance")
				return
			}

			var currentBalance float64
			switch req.Currency {
			case "USD":
				currentBalance = utils.ParseStringToFloat64(wallet.BalanceUsd)
			case "RUB":
				currentBalance = utils.ParseStringToFloat64(wallet.BalanceRub)
			case "EUR":
				currentBalance = utils.ParseStringToFloat64(wallet.BalanceEur)
			}

			if currentBalance < req.Amount {
				utils.RespondWithError(w, http.StatusBadRequest, "Insufficient funds")
				return
			}
		}

		amountStr := fmt.Sprintf("%.2f", operationAmount)

		var err error
		switch req.Currency {
		case "USD":
			err = db.UpdateUSDBalance(r.Context(), database.UpdateUSDBalanceParams{
				Amount: amountStr,
				UserID: userID,
			})
		case "RUB":
			err = db.UpdateRUBBalance(r.Context(), database.UpdateRUBBalanceParams{
				Amount: amountStr,
				UserID: userID,
			})
		case "EUR":
			err = db.UpdateEURBalance(r.Context(), database.UpdateEURBalanceParams{
				Amount: amountStr,
				UserID: userID,
			})
		}

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
			"message": fmt.Sprintf("%s successful", operation),
			"new_balance": map[string]float64{
				"USD": utils.ParseStringToFloat64(updatedWallet.BalanceUsd),
				"RUB": utils.ParseStringToFloat64(updatedWallet.BalanceRub),
				"EUR": utils.ParseStringToFloat64(updatedWallet.BalanceEur),
			},
		})
	}
}

// @Summary Get wallet balance
// @Description Get current balance of user's wallet in all currencies
// @Tags wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.BalanceResponseSwag "Wallet balances"
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

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"balances": map[string]float64{
				"USD": utils.ParseStringToFloat64(wallet.BalanceUsd),
				"RUB": utils.ParseStringToFloat64(wallet.BalanceRub),
				"EUR": utils.ParseStringToFloat64(wallet.BalanceEur),
			},
		})
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
