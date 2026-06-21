package main

import (
	"fmt"
	"github.com/ZhuJincheng-git/stride-backend/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}
