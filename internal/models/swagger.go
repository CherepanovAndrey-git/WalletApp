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
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"secretpass123"`
}

// WalletOperationRequestSwag represents deposit/withdrawal request body
type WalletOperationRequestSwag struct {
	Amount   float64 `json:"amount" example:"100.50"`
	Currency string  `json:"currency" example:"USD"`
}

// BalanceResponseSwag represents the balance response
type BalanceResponseSwag struct {
	Balance float64 `json:"balance" example:"1000.50"`
}
