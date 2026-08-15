package main

import (
	"fmt"

	"inventory/internal/config"
	"inventory/internal/service"
	"inventory/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	svc := service.New(st, cfg)
	_ = svc
	fmt.Println("inventory ready")
}
