package main

import (
	"fmt"
	"github.com/ZhuJincheng-git/stride-backend/internal/config"
)

func main() {
	cfg, _ := config.Load()
	fmt.Printf("Server will run on port: %s in %s mode\n", cfg.ServerPort, cfg.ServerMode)
	fmt.Printf("Database configuration: host=%s, port=%s, user=%s, name=%s\n",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)
}
