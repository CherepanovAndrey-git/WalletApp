package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
	"wallet-app/internal/exchange"

	httpSwagger "github.com/swaggo/http-swagger"
	"wallet-app/internal/database"
	"wallet-app/internal/handlers"
	"wallet-app/internal/middleware"
	"wallet-app/internal/models/apicfg"
	_ "wallet-app/swagger/docs"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load the .env file.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment.")
	}
	fmt.Println("PORT:", port)

	// GRPC
	exchangeClient, err := exchange.NewClient(os.Getenv("EXCHANGER_GRPC_ADDR"))
	if err != nil {
		log.Fatal("Failed to create exchange client:", err)
	}
	defer exchangeClient.Close()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment.")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to the database:", err)
	}

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err = conn.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Successfully connected to database")

	apiCfg := apicfg.ApiConfig{
		DB:     database.New(conn),
		Client: exchangeClient,
	}

	router := chi.NewRouter()
	exchangeHandler := handlers.NewExchangeHandler(apiCfg.Client)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)

	authHandler := handlers.AuthHandler{DB: apiCfg.DB}

	v1Router.Post("/register", authHandler.Register)
	v1Router.Post("/login", authHandler.Login)

	v1Router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/create-wallet", handlers.CreateWalletHandler(apiCfg.DB))
		r.Post("/wallet/deposit", handlers.WalletOperationHandler(apiCfg.DB))
		r.Post("/wallet/withdraw", handlers.WalletOperationHandler(apiCfg.DB))
		r.Get("/balance", handlers.GetBalanceHandler(apiCfg.DB))
		r.Get("/exchange/rates", exchangeHandler.GetRates)
		r.Post("/exchange", exchangeHandler.ExchangeCurrency(apiCfg.DB))
	})

	v1Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/v1/swagger/doc.json"),
	))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Server starting on port %v", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
