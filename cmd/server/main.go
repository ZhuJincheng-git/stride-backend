package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ZhuJincheng-git/stride-backend/internal/app"
	"github.com/ZhuJincheng-git/stride-backend/internal/config"
	"github.com/ZhuJincheng-git/stride-backend/internal/database"
	"github.com/ZhuJincheng-git/stride-backend/internal/model"
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

	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	application := app.New(cfg, db)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.AppPort),
		Handler:           application.Engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%d (mode=%s)", cfg.AppPort, cfg.AppEnv)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
