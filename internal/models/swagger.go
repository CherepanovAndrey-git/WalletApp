package models

// @title Wallet API
// @version 1.0
// @description A simple wallet service API
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// RegisterRequestSwag represents the registration request body
type RegisterRequestSwag struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secretpass123"`
}

// LoginRequestSwag represents the login request body
type LoginRequestSwag struct {
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"secretpass123"`
}

// WalletOperationRequestSwag represents deposit/withdrawal request body
type WalletOperationRequestSwag struct {
	Amount   float64 `json:"amount" example:"100.50"`
	Currency string  `json:"currency" example:"USD"` // Fixed space
}

// BalanceResponseSwag represents the balance response
type BalanceResponseSwag struct {
	Balance float64 `json:"balance" example:"1000.50"`
}

// AuthResponseSwag represents authentication response
type AuthResponseSwag struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ExchangeResponseSwag represents exchange operation response
type ExchangeResponseSwag struct {
	Message         string            `json:"message" example:"Exchange successful"`
	ExchangedAmount string            `json:"exchanged_amount" example:"100.50"`
	NewBalance      map[string]string `json:"new_balance" example:"USD:1500.00,EUR:500.00,RUB:30000.00"`
}
