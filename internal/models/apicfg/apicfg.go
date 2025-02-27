package apicfg

import (
	"wallet-app/internal/database"
	"wallet-app/internal/exchange"
)

type ApiConfig struct {
	DB     *database.Queries
	Client *exchange.Client
}
