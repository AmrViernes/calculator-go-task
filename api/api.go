package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"pack-calculator/calculator"

	"github.com/gorilla/mux"
)

// Server holds our calculator and HTTP router
type Server struct {
	calculator *calculator.Calculator
	router     *mux.Router
}

// NewServer creates a new API server with the given calculator
func NewServer(calc *calculator.Calculator) *Server {
	s := &Server{
		calculator: calc,
		router:     mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

// setupRoutes defines all our API endpoints
func (s *Server) setupRoutes() {
	// Main API endpoints
	s.router.HandleFunc("/api/calculate", s.handleCalculate).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/api/packsizes", s.handleGetPackSizes).Methods("GET")
	s.router.HandleFunc("/api/packsizes", s.handleUpdatePackSizes).Methods("PUT")

	// Serve the frontend UI
	// Try local path first, fallback to container path
	uiPath := "./ui/"
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		uiPath = "/root/ui/"
	}
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir(uiPath)))

	// Enable CORS so frontend can talk to backend
	s.router.Use(corsMiddleware)
}

// CORS middleware - allows cross-origin requests from the frontend
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Pre-flight requests just get a 200 OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetRouter returns the HTTP router (used for testing)
func (s *Server) GetRouter() http.Handler {
	return s.router
}

// Request/response types for JSON parsing
type calculateRequest struct {
	OrderQuantity *int `json:"orderQuantity"`
}

type calculateResponse struct {
	Packs []calculator.PackResult `json:"packs"`
}

// handleCalculate processes order quantity and returns optimal pack breakdown
func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Make sure we actually got an order quantity
	if req.OrderQuantity == nil {
		http.Error(w, "orderQuantity is required", http.StatusBadRequest)
		return
	}

	orderQuantity := *req.OrderQuantity
	if orderQuantity < 0 {
		http.Error(w, "Order quantity must be non-negative", http.StatusBadRequest)
		return
	}

	// Calculate optimal packs and send back the result
	packs := s.calculator.Calculate(orderQuantity)
	json.NewEncoder(w).Encode(calculateResponse{Packs: packs})
}

// handleGetPackSizes returns the current pack sizes in use
func (s *Server) handleGetPackSizes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sizes := s.calculator.GetPackSizes()

	type packSizeResponse struct {
		PackSizes []int `json:"packSizes"`
	}
	json.NewEncoder(w).Encode(packSizeResponse{PackSizes: sizes})
}

type updatePackSizesRequest struct {
	PackSizes []int `json:"packSizes"`
}

// handleUpdatePackSizes updates the available pack sizes
func (s *Server) handleUpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req updatePackSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate we have at least one pack size
	if len(req.PackSizes) == 0 {
		http.Error(w, "At least one pack size is required", http.StatusBadRequest)
		return
	}

	// All pack sizes must be positive numbers
	for _, size := range req.PackSizes {
		if size <= 0 {
			http.Error(w, "All pack sizes must be positive", http.StatusBadRequest)
			return
		}
	}

	// Update and echo back the new sizes
	s.calculator.UpdatePackSizes(req.PackSizes)

	type packSizeResponse struct {
		PackSizes []int `json:"packSizes"`
	}
	json.NewEncoder(w).Encode(packSizeResponse{PackSizes: s.calculator.GetPackSizes()})
}

// Start begins listening for HTTP requests
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.GetRouter())
}
