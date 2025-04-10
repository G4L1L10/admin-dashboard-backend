package main

import (
	"log"

	"github.com/G4L1L10/admin-dashboard-backend/internal/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
