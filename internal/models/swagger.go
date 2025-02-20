package models

// @title Wallet API
// @version 1.0
// @description A simple wallet service API
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// RegisterRequest represents the registration request body
type RegisterRequestSwag struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secretpass123"`
}

// LoginRequest represents the login request body
type LoginRequestSwag struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secretpass123"`
}

// WalletOperationRequest represents deposit/withdrawal request body
type WalletOperationRequestSwag struct {
	Amount   float64 `json:"amount" example:"100.50"`
	Currency string  `json:"currency" example:"USD"`
}

// BalanceResponse represents the balance response
type BalanceResponseSwag struct {
	Balance float64 `json:"balance" example:"1000.50"`
}
