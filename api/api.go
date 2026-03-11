package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"pack-calculator/calculator"

	"github.com/gorilla/mux"
)

type Server struct {
	calculator *calculator.Calculator
	router     *mux.Router
}

func NewServer(calc *calculator.Calculator) *Server {
	s := &Server{
		calculator: calc,
		router:     mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/api/calculate", s.handleCalculate).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/api/packsizes", s.handleGetPackSizes).Methods("GET")
	s.router.HandleFunc("/api/packsizes", s.handleUpdatePackSizes).Methods("PUT")

	uiPath := "./ui/"
	if _, err := os.Stat(uiPath); os.IsNotExist(err) {
		uiPath = "/root/ui/"
	}
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir(uiPath)))

	s.router.Use(corsMiddleware)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) GetRouter() http.Handler {
	return s.router
}

type calculateRequest struct {
	OrderQuantity *int `json:"orderQuantity"`
}

type calculateResponse struct {
	Packs []calculator.PackResult `json:"packs"`
}

func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.OrderQuantity == nil {
		http.Error(w, "orderQuantity is required", http.StatusBadRequest)
		return
	}

	orderQuantity := *req.OrderQuantity
	if orderQuantity < 0 {
		http.Error(w, "Order quantity must be non-negative", http.StatusBadRequest)
		return
	}

	packs := s.calculator.Calculate(orderQuantity)

	json.NewEncoder(w).Encode(calculateResponse{Packs: packs})
}

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

func (s *Server) handleUpdatePackSizes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req updatePackSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.PackSizes) == 0 {
		http.Error(w, "At least one pack size is required", http.StatusBadRequest)
		return
	}

	for _, size := range req.PackSizes {
		if size <= 0 {
			http.Error(w, "All pack sizes must be positive", http.StatusBadRequest)
			return
		}
	}

	s.calculator.UpdatePackSizes(req.PackSizes)

	type packSizeResponse struct {
		PackSizes []int `json:"packSizes"`
	}
	json.NewEncoder(w).Encode(packSizeResponse{PackSizes: s.calculator.GetPackSizes()})
}

func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.GetRouter())
}
