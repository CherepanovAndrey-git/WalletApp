package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

// WalletOperationRequest - структура для депозита и снятия средств.
type WalletOperationRequest struct {
	WalletID      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}

// WalletOperationHandler - функция для типов операций над кошельком, DEPOSIT или WITHDRAWAL

// @Summary Perform a wallet operation (DEPOSIT or WITHDRAW)
// @Description Deposit or withdraw funds from a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param ApiKey header string true "API Key"
// @Param operation body WalletOperationRequest true "Wallet operation details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /wallet [post]
func WalletOperationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверка API Key
		apiKey := r.Header.Get("ApiKey")
		if apiKey == "" {
			http.Error(w, "ApiKey header missing", http.StatusUnauthorized)
			return
		}

		//
		userID, err := ValidateAPIKey(db, apiKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
			} else {
				log.Printf("Error validating API key: %v", err)
				http.Error(w, "Failed to validate API key", http.StatusInternalServerError)
			}
			return
		}

		// Парсинг тела запроса
		var req WalletOperationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Валидация суммы
		if req.Amount <= 0 {
			http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
			return
		}

		// Выбор операции (DEPOSIT или WITHDRAW)
		var query string
		switch req.OperationType {
		case "DEPOSIT":
			query = `UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE uuid = $2 AND user_id = $3 RETURNING balance`
		case "WITHDRAW":
			query = `UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE uuid = $2 AND user_id = $3 AND balance >= $1 RETURNING balance`
		default:
			http.Error(w, "Invalid operation type", http.StatusBadRequest)
			return
		}

		// Выполнение операции
		var newBalance float64
		err = db.QueryRow(query, req.Amount, req.WalletID, userID).Scan(&newBalance)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				if req.OperationType == "WITHDRAW" {
					http.Error(w, "Insufficient balance or wallet not found", http.StatusBadRequest)
				} else {
					http.Error(w, "Wallet not found", http.StatusNotFound)
				}
			} else {
				log.Printf("Error updating balance: %v", err)
				http.Error(w, "Failed to update balance", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"walletId": req.WalletID,
			"balance":  newBalance,
		})
	}
}

func GetBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("ApiKey")
		if apiKey == "" {
			http.Error(w, "ApiKey header missing", http.StatusUnauthorized)
			return
		}

		// Валидация АПИ.
		userID, err := ValidateAPIKey(db, apiKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
			} else {
				log.Printf("Error validating API key: %v", err)
				http.Error(w, "Failed to validate API key", http.StatusInternalServerError)
			}
			return
		}
		log.Printf("API key validated for user ID: %s", userID)

		// извлекаем uuid с URL используя Chi
		walletUUID := chi.URLParam(r, "uuid")
		log.Printf("Extracted wallet UUID: %s", walletUUID)

		if walletUUID == "" {
			http.Error(w, "Wallet UUID is required", http.StatusBadRequest)
			return
		}

		// Запрашивает баланс с БД
		var balance float64
		err = db.QueryRow("SELECT balance FROM wallets WHERE uuid = $1 AND user_id = $2", walletUUID, userID).Scan(&balance)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Wallet not found for UUID: %s and user ID: %s", walletUUID, userID)
				http.Error(w, "Wallet not found", http.StatusNotFound)
			} else {
				log.Printf("Error fetching balance: %v", err)
				http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
			}
			return
		}

		// Возвращает баланс
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{
			"balance": balance,
		})
	}
}
