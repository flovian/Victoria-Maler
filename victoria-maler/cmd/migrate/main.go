package main

import (
	"log"

	"ecochain-victoria/internal/config"
	"ecochain-victoria/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db, "migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}
	log.Println("database migrations applied successfully")
}
