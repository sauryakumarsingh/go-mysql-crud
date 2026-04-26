package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourname/go-mysql-crud/internal/db"
	"github.com/yourname/go-mysql-crud/internal/router"
)

func main() {
	// Load .env (ignored if not present — production uses real env vars)
	_ = godotenv.Load()

	// Connect to MySQL
	database, err := db.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Could not connect to MySQL: %v\n", err)
		fmt.Println("   Make sure MySQL is running and .env is configured correctly.")
		os.Exit(1)
	}
	defer database.Close()
	fmt.Println("✅ Connected to MySQL")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router.New(database),
	}

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
