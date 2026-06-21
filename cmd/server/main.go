package main

import (
	"fmt"
	"log"

	"github.com/ZhuJincheng-git/stride-backend/internal/config"
	"github.com/ZhuJincheng-git/stride-backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	_, err = database.Open(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
}
