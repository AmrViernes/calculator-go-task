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
	port := 8080

	if portStr := os.Getenv("PORT"); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	calc := calculator.NewDefaultCalculator()

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

	server := api.NewServer(calc)

	log.Printf("Starting Pack Calculator API")
	log.Printf("Default pack sizes: %v", calc.GetPackSizes())
	log.Printf("Server starting on http://localhost:%d", port)

	if err := server.Start(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
