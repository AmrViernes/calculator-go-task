package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"pack-calculator/api"
	"pack-calculator/calculator"
)

func main() {
	// Default port, but can be overridden via environment variable
	port := 8080

	// Check if PORT env var is set (useful for cloud deployments)
	if portStr := os.Getenv("PORT"); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	// Create calculator with default pack sizes
	calc := calculator.NewDefaultCalculator()

	// Allow custom pack sizes via environment variable
	// Format: "250,500,1000,2000,5000"
	if packSizesEnv := os.Getenv("PACK_SIZES"); packSizesEnv != "" {
		var sizes []int
		for _, s := range strings.Split(packSizesEnv, ",") {
			var size int
			s = strings.TrimSpace(s)
			if _, err := fmt.Sscanf(s, "%d", &size); err == nil && size > 0 {
				sizes = append(sizes, size)
			}
		}
		if len(sizes) > 0 {
			calc.UpdatePackSizes(sizes)
			log.Printf("Using custom pack sizes from environment: %v", sizes)
		}
	}

	// Wire up the API server with our calculator
	server := api.NewServer(calc)

	log.Printf("Starting Pack Calculator API")
	log.Printf("Default pack sizes: %v", calc.GetPackSizes())
	log.Printf("Server starting on http://localhost:%d", port)

	if err := server.Start(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
