package main

import (
	"context"
	"log"

	"go-srv-temp/internal/app"
	"go-srv-temp/internal/config"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	application, err := app.New(context.Background(), cfg)
	if err != nil {
		log.Fatalf("create app: %v", err)
	}

	if err := application.Run(context.Background()); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
