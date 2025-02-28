package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ParseStringToFloat64 converts a string to float64.
func ParseStringToFloat64(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		log.Printf("Failed to parse balance: %v", err)
		return 0
	}
	return parsed
}

// RespondWithError sends an error response.
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response.
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
