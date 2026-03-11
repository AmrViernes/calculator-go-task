package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pack-calculator/calculator"
)

// TestHandleCalculate tests the calculate endpoint
func TestHandleCalculate(t *testing.T) {
	calc := calculator.NewDefaultCalculator()
	server := NewServer(calc)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		validateFunc   func(t *testing.T, body string)
	}{
		{
			name:           "valid request - 1 item",
			requestBody:    `{"orderQuantity": 1}`,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body string) {
				var resp calculateResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if len(resp.Packs) != 1 || resp.Packs[0].Size != 250 || resp.Packs[0].Count != 1 {
					t.Errorf("Expected 1x250, got %v", resp.Packs)
				}
			},
		},
		{
			name:           "valid request - 251 items",
			requestBody:    `{"orderQuantity": 251}`,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body string) {
				var resp calculateResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if len(resp.Packs) != 1 || resp.Packs[0].Size != 500 || resp.Packs[0].Count != 1 {
					t.Errorf("Expected 1x500, got %v", resp.Packs)
				}
			},
		},
		{
			name:           "valid request - 501 items",
			requestBody:    `{"orderQuantity": 501}`,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body string) {
				var resp calculateResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if len(resp.Packs) != 2 {
					t.Errorf("Expected 2 pack types, got %d: %v", len(resp.Packs), resp.Packs)
				}
			},
		},
		{
			name:           "valid request - 0 items",
			requestBody:    `{"orderQuantity": 0}`,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body string) {
				var resp calculateResponse
				if err := json.Unmarshal([]byte(body), &resp); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if len(resp.Packs) != 0 {
					t.Errorf("Expected 0 packs, got %v", resp.Packs)
				}
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"orderQuantity": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			},
		{
			name:           "negative quantity",
			requestBody:    `{"orderQuantity": -1}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing field",
			requestBody:    `{}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/calculate", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleCalculate(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, w.Body.String())
			}
		})
	}
}

// TestHandleGetPackSizes tests the get pack sizes endpoint
func TestHandleGetPackSizes(t *testing.T) {
	calc := calculator.NewDefaultCalculator()
	server := NewServer(calc)

	req := httptest.NewRequest("GET", "/api/packsizes", nil)
	w := httptest.NewRecorder()

	server.handleGetPackSizes(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp struct {
		PackSizes []int `json:"packSizes"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	expected := []int{250, 500, 1000, 2000, 5000}
	if len(resp.PackSizes) != len(expected) {
		t.Errorf("Expected %d pack sizes, got %d", len(expected), len(resp.PackSizes))
	}

	for i, size := range resp.PackSizes {
		if size != expected[i] {
			t.Errorf("Expected pack size %d at position %d, got %d", expected[i], i, size)
		}
	}
}

// TestHandleUpdatePackSizes tests the update pack sizes endpoint
func TestHandleUpdatePackSizes(t *testing.T) {
	calc := calculator.NewDefaultCalculator()
	server := NewServer(calc)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		validateFunc   func(t *testing.T, body string, calc *calculator.Calculator)
	}{
		{
			name:           "valid update",
			requestBody:    `{"packSizes": [100, 200, 400]}`,
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, body string, calc *calculator.Calculator) {
				sizes := calc.GetPackSizes()
				expected := []int{100, 200, 400}
				if len(sizes) != len(expected) {
					t.Errorf("Expected %d pack sizes, got %d", len(expected), len(sizes))
				}
			},
		},
		{
			name:           "empty pack sizes",
			requestBody:    `{"packSizes": []}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative pack size",
			requestBody:    `{"packSizes": [100, -200, 300]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "zero pack size",
			requestBody:    `{"packSizes": [100, 0, 300]}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"packSizes": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset calculator before each test
			calc = calculator.NewDefaultCalculator()
			server = NewServer(calc)

			req := httptest.NewRequest("PUT", "/api/packsizes", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleUpdatePackSizes(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, w.Body.String(), calc)
			}
		})
	}
}

// TestCORSMiddleware tests CORS headers
func TestCORSMiddleware(t *testing.T) {
	calc := calculator.NewDefaultCalculator()
	server := NewServer(calc)

	tests := []struct {
		name           string
		method         string
		expectedOrigin string
	}{
		{"OPTIONS request", "OPTIONS", "*"},
		{"GET request", "GET", "*"},
		{"POST request", "POST", "*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/calculate", nil)
			w := httptest.NewRecorder()

			server.GetRouter().ServeHTTP(w, req)

			origin := w.Header().Get("Access-Control-Allow-Origin")
			if origin != tt.expectedOrigin {
				t.Errorf("Expected CORS origin %s, got %s", tt.expectedOrigin, origin)
			}
		})
	}
}

// TestCalculateIntegration tests full calculate integration
func TestCalculateIntegration(t *testing.T) {
	calc := calculator.NewDefaultCalculator()
	server := NewServer(calc)

	// Test calculate, update pack sizes, then calculate again
	testCases := []struct {
		name          string
		packSizes     []int
		orderQuantity int
		expectedPacks []calculator.PackResult
	}{
		{
			name:          "default packs - 251",
			packSizes:     nil, // use default
			orderQuantity: 251,
			expectedPacks: []calculator.PackResult{{Size: 500, Count: 1}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.packSizes != nil {
				// Update pack sizes
				body, _ := json.Marshal(map[string]interface{}{
					"packSizes": tt.packSizes,
				})
				req := httptest.NewRequest("PUT", "/api/packsizes", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				server.handleUpdatePackSizes(w, req)

				if w.Code != http.StatusOK {
					t.Errorf("Failed to update pack sizes: %d", w.Code)
				}
			}

			// Calculate
			body, _ := json.Marshal(map[string]interface{}{
				"orderQuantity": tt.orderQuantity,
			})
			req := httptest.NewRequest("POST", "/api/calculate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleCalculate(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}

			var resp calculateResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if len(resp.Packs) != len(tt.expectedPacks) {
				t.Errorf("Expected %d pack types, got %d: %v", len(tt.expectedPacks), len(resp.Packs), resp.Packs)
			}

			for i, pack := range tt.expectedPacks {
				if resp.Packs[i].Size != pack.Size || resp.Packs[i].Count != pack.Count {
					t.Errorf("Expected pack %+v at position %d, got %+v", pack, i, resp.Packs[i])
				}
			}
		})
	}
}
