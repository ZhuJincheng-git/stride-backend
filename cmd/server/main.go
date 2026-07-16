package main

import (
	"log"

	"github.com/ZhuJincheng-git/stride-backend/internal/config"
	"github.com/ZhuJincheng-git/stride-backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("db handle: %v", err)
	}
	defer sqlDB.Close()
}
