package main

import (
	"log"
	"os"

	"payment/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Overload(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("gagal membaca file .env: %v", err)
	}

	config := database.ConfigFromEnv()
	db, err := database.Open(config)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("gagal menutup koneksi database: %v", err)
		}
	}()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("gagal menjalankan migration: %v", err)
	}

	log.Println("migration berhasil dijalankan")
}
