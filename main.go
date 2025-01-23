package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"main.go/internal/database"
	"net/http"
	"os"
	"time"
	_ "github.com/lib/pq"
	"github.com/swaggo/http-swagger"
	_"main.go/docs"
)

type apiConfig struct {
	DB *database.Queries
}

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

	
	if err := conn.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Successfully connected to database")

	apiCfg := apiConfig{DB: database.New(conn)}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", HandlerReadiness)
	v1Router.Get("/err", HandlerErr)
	v1Router.Post("/users", apiCfg.HandlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/wallet", WalletOperationHandler(conn))
	v1Router.Get("/wallets/{uuid}/balance", GetBalanceHandler(conn))
	v1Router.Post("/wallets", CreateWalletHandler(conn))
	v1Router.Get("/swagger/*", httpSwagger.WrapHandler) //TODO: Annotations for swagger, this is mock for now.

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
