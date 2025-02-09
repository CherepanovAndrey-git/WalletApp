package apicfg

import "wallet-app/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
