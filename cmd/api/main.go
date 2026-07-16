package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"payment/internal/bootstrap"
	"payment/internal/database"
)

func main() {
	if err := godotenv.Overload(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("gagal membaca file .env: %v", err)
	}

	dbConfig := database.ConfigFromEnv()

	db, err := database.Open(dbConfig)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("gagal menjalankan auto migration: %v", err)
	}

	app := bootstrap.NewApp(db)
	port := getEnv("APP_PORT", "8080")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("server berjalan pada http://localhost:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gagal menjalankan server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("gagal menghentikan server dengan baik: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
