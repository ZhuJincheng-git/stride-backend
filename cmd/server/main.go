package main

import (
	"fmt"
	"github.com/ZhuJincheng-git/stride-backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
