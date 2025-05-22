package main

import (
	"MyStonks-go/internal/config"
	"MyStonks-go/internal/server"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize and start server
	srv := server.NewServer(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
