package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
	"time"
	"wallet-app/internal/database"
	"wallet-app/internal/exchange"
	"wallet-app/internal/utils"
)

type ExchangeRequest struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
}

type ExchangeHandler struct {
	client *exchange.Client
	cache  struct {
		sync.RWMutex
		rates       map[string]float32
		lastUpdated time.Time
	}
}

func NewExchangeHandler(client *exchange.Client) *ExchangeHandler {
	return &ExchangeHandler{
		client: client,
		cache: struct {
			sync.RWMutex
			rates       map[string]float32
			lastUpdated time.Time
		}{
			rates: make(map[string]float32),
		},
	}
}

func (h *ExchangeHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	log.Println("HTTP GetRates endpoint called")

	rates, err := h.client.GetRates(r.Context())
	if err != nil {
		log.Printf("gRPC error: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get rates")
		return
	}
	log.Printf("Received %d rates from gRPC", len(rates))
	respondWithFormattedRates(w, rates)
}

func (h *ExchangeHandler) ExchangeCurrency(db *database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		var req ExchangeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if !isValidCurrency(req.FromCurrency) || !isValidCurrency(req.ToCurrency) {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid currency pair")
			return
		}

		rate, err := h.getExchangeRate(req.FromCurrency, req.ToCurrency)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get exchange rate")
			return
		}

		log.Printf("Attempting to get wallet for user: %v", userID)
		wallet, err := db.GetWalletByUserID(r.Context(), userID)
		if err != nil {
			log.Printf("Database error: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get wallet")
			return
		}

		fromBalance := getCurrencyBalance(wallet, req.FromCurrency)
		if fromBalance < req.Amount {
			utils.RespondWithError(w, http.StatusBadRequest, "Insufficient funds")
			return
		}

		exchangedAmount := req.Amount * float64(rate)
		toBalance := getCurrencyBalance(wallet, req.ToCurrency)

		fromDelta := -req.Amount
		toDelta := exchangedAmount

		err = updateBalances(db, r.Context(), userID, req.FromCurrency, fromDelta, req.ToCurrency, toDelta)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update balances")
			return
		}

		newFromBalance := fromBalance + fromDelta
		newToBalance := toBalance + toDelta

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"message":          "Exchange successful",
			"exchanged_amount": fmt.Sprintf("%.2f", exchangedAmount),
			"new_balance": map[string]string{
				req.FromCurrency: fmt.Sprintf("%.2f", newFromBalance),
				req.ToCurrency:   fmt.Sprintf("%.2f", newToBalance),
			},
		})
		log.Printf("User ID: %s", userID.String())
		log.Printf("From Currency: %s, To Currency: %s, Amount: %.2f", req.FromCurrency, req.ToCurrency, req.Amount)
		log.Printf("From Balance: %.2f", fromBalance)
		log.Printf("Exchanged Amount: %.2f", exchangedAmount)
		log.Printf("New From Balance: %.2f, New To Balance: %.2f", newFromBalance, newToBalance)
	}
}

func (h *ExchangeHandler) getExchangeRate(from, to string) (float32, error) {
	h.cache.RLock()
	defer h.cache.RUnlock()

	key := fmt.Sprintf("%s_%s", from, to)
	if rate, exists := h.cache.rates[key]; exists {
		return rate, nil
	}

	return h.client.GetRate(from, to)
}

func respondWithFormattedRates(w http.ResponseWriter, rates map[string]float32) {
	formatted := make(map[string]float32)
	for pair, rate := range rates {
		formatted[pair] = rate
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"rates": formatted})
}

func isValidCurrency(currency string) bool {
	valid := map[string]bool{"USD": true, "EUR": true, "RUB": true}
	return valid[currency]
}

func getCurrencyBalance(wallet database.Wallet, currency string) float64 {
	switch currency {
	case "USD":
		return utils.ParseStringToFloat64(wallet.BalanceUsd)
	case "EUR":
		return utils.ParseStringToFloat64(wallet.BalanceEur)
	case "RUB":
		return utils.ParseStringToFloat64(wallet.BalanceRub)
	default:
		return 0
	}
}

func updateBalances(db *database.Queries, ctx context.Context, userID uuid.UUID,
	fromCurrency string, fromDelta float64, toCurrency string, toDelta float64) error {

	fromUpdate := fmt.Sprintf("%.2f", fromDelta)
	switch fromCurrency {
	case "USD":
		err := db.UpdateUSDBalance(ctx, database.UpdateUSDBalanceParams{
			UserID: userID,
			Amount: fromUpdate,
		})
		if err != nil {
			return err
		}
	case "EUR":
		err := db.UpdateEURBalance(ctx, database.UpdateEURBalanceParams{
			UserID: userID,
			Amount: fromUpdate,
		})
		if err != nil {
			return err
		}
	case "RUB":
		err := db.UpdateRUBBalance(ctx, database.UpdateRUBBalanceParams{
			UserID: userID,
			Amount: fromUpdate,
		})
		if err != nil {
			return err
		}
	}

	toUpdate := fmt.Sprintf("%.2f", toDelta)
	switch toCurrency {
	case "USD":
		err := db.UpdateUSDBalance(ctx, database.UpdateUSDBalanceParams{
			UserID: userID,
			Amount: toUpdate,
		})
		if err != nil {
			return err
		}
	case "EUR":
		err := db.UpdateEURBalance(ctx, database.UpdateEURBalanceParams{
			UserID: userID,
			Amount: toUpdate,
		})
		if err != nil {
			return err
		}
	case "RUB":
		err := db.UpdateRUBBalance(ctx, database.UpdateRUBBalanceParams{
			UserID: userID,
			Amount: toUpdate,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
