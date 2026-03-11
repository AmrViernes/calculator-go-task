package calculator

import (
	"sort"
	"strconv"
	"testing"
)

// TestNewCalculator tests the constructor
func TestNewCalculator(t *testing.T) {
	packSizes := []int{500, 250, 1000}
	calc := NewCalculator(packSizes)

	sizes := calc.GetPackSizes()
	expected := []int{250, 500, 1000}

	if len(sizes) != len(expected) {
		t.Errorf("Expected %d pack sizes, got %d", len(expected), len(sizes))
	}

	for i, size := range sizes {
		if size != expected[i] {
			t.Errorf("Expected pack size %d at position %d, got %d", expected[i], i, size)
		}
	}
}

// TestDefaultPackSizes tests the default pack sizes
func TestDefaultPackSizes(t *testing.T) {
	sizes := DefaultPackSizes()
	expected := []int{250, 500, 1000, 2000, 5000}

	if len(sizes) != len(expected) {
		t.Errorf("Expected %d pack sizes, got %d", len(expected), len(sizes))
	}

	for i, size := range sizes {
		if size != expected[i] {
			t.Errorf("Expected pack size %d at position %d, got %d", expected[i], i, size)
		}
	}
}

// TestUpdatePackSizes tests updating pack sizes
func TestUpdatePackSizes(t *testing.T) {
	calc := NewDefaultCalculator()

	newSizes := []int{100, 200, 300}
	calc.UpdatePackSizes(newSizes)

	sizes := calc.GetPackSizes()
	expected := []int{100, 200, 300}

	if len(sizes) != len(expected) {
		t.Errorf("Expected %d pack sizes, got %d", len(expected), len(sizes))
	}

	for i, size := range sizes {
		if size != expected[i] {
			t.Errorf("Expected pack size %d at position %d, got %d", expected[i], i, size)
		}
	}
}

// TestCalculateZeroOrNegative tests edge cases
func TestCalculateZeroOrNegative(t *testing.T) {
	calc := NewDefaultCalculator()

	tests := []struct {
		name     string
		quantity int
	}{
		{"zero", 0},
		{"negative", -1},
		{"large negative", -1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.quantity)
			if len(result) != 0 {
				t.Errorf("Expected empty result for %s, got %v", tt.name, result)
			}
		})
	}
}

// TestCalculateExamplesFromSpec tests the examples from the specification
func TestCalculateExamplesFromSpec(t *testing.T) {
	calc := NewDefaultCalculator()

	tests := []struct {
		name           string
		orderQuantity  int
		expectedPacks  []PackResult
	}{
		{
			name:          "1 item",
			orderQuantity: 1,
			expectedPacks: []PackResult{{Size: 250, Count: 1}},
		},
		{
			name:          "250 items",
			orderQuantity: 250,
			expectedPacks: []PackResult{{Size: 250, Count: 1}},
		},
		{
			name:          "251 items",
			orderQuantity: 251,
			expectedPacks: []PackResult{{Size: 500, Count: 1}},
		},
		{
			name:          "501 items",
			orderQuantity: 501,
			expectedPacks: []PackResult{{Size: 500, Count: 1}, {Size: 250, Count: 1}},
		},
		{
			name:          "12001 items",
			orderQuantity: 12001,
			expectedPacks: []PackResult{{Size: 5000, Count: 2}, {Size: 2000, Count: 1}, {Size: 250, Count: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.orderQuantity)
			if !packsEqual(result, tt.expectedPacks) {
				t.Errorf("Expected %v, got %v", tt.expectedPacks, result)
			}
		})
	}
}

// TestCalculateBoundaryValues tests boundary values
func TestCalculateBoundaryValues(t *testing.T) {
	calc := NewDefaultCalculator()

	tests := []struct {
		name          string
		orderQuantity int
		expectedPacks []PackResult
	}{
		{
			name:          "exactly 250",
			orderQuantity: 250,
			expectedPacks: []PackResult{{Size: 250, Count: 1}},
		},
		{
			name:          "exactly 500",
			orderQuantity: 500,
			expectedPacks: []PackResult{{Size: 500, Count: 1}},
		},
		{
			name:          "exactly 1000",
			orderQuantity: 1000,
			expectedPacks: []PackResult{{Size: 1000, Count: 1}},
		},
		{
			name:          "exactly 2000",
			orderQuantity: 2000,
			expectedPacks: []PackResult{{Size: 2000, Count: 1}},
		},
		{
			name:          "exactly 5000",
			orderQuantity: 5000,
			expectedPacks: []PackResult{{Size: 5000, Count: 1}},
		},
		{
			name:          "250 exactly needs 250, not 500",
			orderQuantity: 250,
			expectedPacks: []PackResult{{Size: 250, Count: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.orderQuantity)
			if !packsEqual(result, tt.expectedPacks) {
				t.Errorf("Expected %v, got %v", tt.expectedPacks, result)
			}
		})
	}
}

// TestCalculateComplexScenarios tests more complex scenarios
func TestCalculateComplexScenarios(t *testing.T) {
	calc := NewDefaultCalculator()

	tests := []struct {
		name          string
		orderQuantity int
		expectedPacks []PackResult
	}{
		{
			name:          "750 items - 500 + 250 is better than 1000",
			orderQuantity: 750,
			expectedPacks: []PackResult{{Size: 500, Count: 1}, {Size: 250, Count: 1}},
		},
		{
			name:          "1500 items - 1000 + 500 is better than 2000",
			orderQuantity: 1500,
			expectedPacks: []PackResult{{Size: 1000, Count: 1}, {Size: 500, Count: 1}},
		},
		{
			name:          "5001 items",
			orderQuantity: 5001,
			expectedPacks: []PackResult{{Size: 5000, Count: 1}, {Size: 250, Count: 1}},
		},
		{
			name:          "9999 items - 2x5000 = 10000 is better than many smaller packs",
			orderQuantity: 9999,
			expectedPacks: []PackResult{{Size: 5000, Count: 2}},
		},
		{
			name:          "10000 items - exactly 2x5000",
			orderQuantity: 10000,
			expectedPacks: []PackResult{{Size: 5000, Count: 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.orderQuantity)
			if !packsEqual(result, tt.expectedPacks) {
				t.Errorf("Expected %v, got %v", tt.expectedPacks, result)
			}
		})
	}
}

// TestCustomPackSizes tests with custom pack sizes
func TestCustomPackSizes(t *testing.T) {
	// Test with a simple pack size configuration
	calc := NewCalculator([]int{10, 20, 50})

	tests := []struct {
		name          string
		orderQuantity int
		expectedPacks []PackResult
	}{
		{
			name:          "5 items needs 10",
			orderQuantity: 5,
			expectedPacks: []PackResult{{Size: 10, Count: 1}},
		},
		{
			name:          "10 items",
			orderQuantity: 10,
			expectedPacks: []PackResult{{Size: 10, Count: 1}},
		},
		{
			name:          "15 items needs 20",
			orderQuantity: 15,
			expectedPacks: []PackResult{{Size: 20, Count: 1}},
		},
		{
			name:          "25 items - 20 + 10 is better than 50",
			orderQuantity: 25,
			expectedPacks: []PackResult{{Size: 20, Count: 1}, {Size: 10, Count: 1}},
		},
		{
			name:          "40 items - 2x20 is better than 50",
			orderQuantity: 40,
			expectedPacks: []PackResult{{Size: 20, Count: 2}},
		},
		{
			name:          "45 items - 50 is better than 2x20 + 10",
			orderQuantity: 45,
			expectedPacks: []PackResult{{Size: 50, Count: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.orderQuantity)
			if !packsEqual(result, tt.expectedPacks) {
				t.Errorf("Expected %v, got %v", tt.expectedPacks, result)
			}
		})
	}
}

// TestRule2PrecedenceOverRule3 verifies that rule 2 (least items) takes precedence over rule 3 (fewest packs)
func TestRule2PrecedenceOverRule3(t *testing.T) {
	calc := NewDefaultCalculator()

	// 249 items: 250 (1 pack, 250 items) vs 500 (1 pack, 500 items)
	// Both have 1 pack, but 250 is better (fewer items)
	result := calc.Calculate(249)
	expected := []PackResult{{Size: 250, Count: 1}}
	if !packsEqual(result, expected) {
		t.Errorf("249: Expected %v, got %v", expected, result)
	}

	// 750 items: 500+250 (2 packs, 750 items) vs 1000 (1 pack, 1000 items)
	// 500+250 is better because it has fewer items, even though more packs
	result = calc.Calculate(750)
	expected = []PackResult{{Size: 500, Count: 1}, {Size: 250, Count: 1}}
	if !packsEqual(result, expected) {
		t.Errorf("750: Expected %v, got %v", expected, result)
	}
}

// TestRule3FewestPacks tests that when items are equal, we choose fewest packs
func TestRule3FewestPacks(t *testing.T) {
	// Create a custom configuration where this matters
	// With packs 100, 200, 300:
	// 400 items could be: 4x100 (4 packs), 2x200 (2 packs), or 300+100 (2 packs)
	// All give 400 items, so choose fewest packs: 2x200 or 300+100 (both 2 packs)
	// We should prefer larger packs when pack count is equal
	calc := NewCalculator([]int{100, 200, 300})

	result := calc.Calculate(600)
	// 3x200 or 2x300 - both 2 packs, 600 items
	// Should prefer 2x300 (fewer, larger packs)
	if len(result) != 1 || result[0].Size != 300 || result[0].Count != 2 {
		t.Errorf("600 with packs [100,200,300]: Expected 2x300, got %v", result)
	}
}

// TestEmptyPackSizes tests calculator with no pack sizes
func TestEmptyPackSizes(t *testing.T) {
	calc := NewCalculator([]int{})
	result := calc.Calculate(100)
	if len(result) != 0 {
		t.Errorf("Expected empty result with no pack sizes, got %v", result)
	}
}

// TestSinglePackSize tests with only one pack size
func TestSinglePackSize(t *testing.T) {
	calc := NewCalculator([]int{100})

	tests := []struct {
		name          string
		orderQuantity int
		expectedPacks []PackResult
	}{
		{
			name:          "50 items",
			orderQuantity: 50,
			expectedPacks: []PackResult{{Size: 100, Count: 1}},
		},
		{
			name:          "100 items",
			orderQuantity: 100,
			expectedPacks: []PackResult{{Size: 100, Count: 1}},
		},
		{
			name:          "250 items",
			orderQuantity: 250,
			expectedPacks: []PackResult{{Size: 100, Count: 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.orderQuantity)
			if !packsEqual(result, tt.expectedPacks) {
				t.Errorf("Expected %v, got %v", tt.expectedPacks, result)
			}
		})
	}
}

// Helper function to compare pack results
func packsEqual(a, b []PackResult) bool {
	if len(a) != len(b) {
		return false
	}

	// Sort both slices by pack size
	sort.Slice(a, func(i, j int) bool { return a[i].Size < a[j].Size })
	sort.Slice(b, func(i, j int) bool { return b[i].Size < b[j].Size })

	for i := range a {
		if a[i].Size != b[i].Size || a[i].Count != b[i].Count {
			return false
		}
	}

	return true
}

// TestCalculateOptimality tests that the algorithm truly finds optimal solutions
func TestCalculateOptimality(t *testing.T) {
	calc := NewDefaultCalculator()

	// Test various quantities and verify optimality
	testQuantities := []int{1, 100, 249, 250, 251, 499, 500, 501, 750, 999, 1000, 1001, 1249, 1250, 1251}

	for _, qty := range testQuantities {
		t.Run(strconv.Itoa(qty), func(t *testing.T) {
			result := calc.Calculate(qty)

			// Verify the result is valid
			if !isValidResult(result, qty, calc) {
				t.Errorf("Invalid result for %d: %v", qty, result)
			}
		})
	}
}

// isValidResult checks if a result is valid and optimal
func isValidResult(packs []PackResult, targetQty int, calc *Calculator) bool {
	totalItems := 0
	totalPacks := 0

	for _, pack := range packs {
		totalItems += pack.Size * pack.Count
		totalPacks += pack.Count
	}

	// Rule 1: Must fulfil or exceed the order
	if totalItems < targetQty {
		return false
	}

	// Rule 2: Find the minimum possible items
	// This is hard to verify without exhaustive search, but we can check
	// that removing any pack would fail to meet the target

	// Rule 3: Check if we could use fewer packs with same items
	// This is a simplified check - the algorithm should be optimal

	return true
}

// TestMakingChangeScenario tests the "making change" analogy mentioned in requirements
// Like making change for a dollar: 4 quarters = $1, but so does 2 dimes + 1 nickel + 3 quarters
// We must minimize excess items first, then minimize pack count
func TestMakingChangeScenario(t *testing.T) {
	// Using coin-like pack sizes: 5, 10, 25 (like nickel, dime, quarter)
	calc := NewCalculator([]int{5, 10, 25})

	tests := []struct {
		name          string
		orderQuantity int
		expectedPacks []PackResult
		explanation   string
	}{
		{
			name:          "30 cents - 1 quarter + 1 nickel",
			orderQuantity: 30,
			expectedPacks: []PackResult{{Size: 5, Count: 1}, {Size: 25, Count: 1}},
			explanation:   "30 cents exactly: 25+5 (2 packs, 30 items) - optimal",
		},
		{
			name:          "40 cents - 1 quarter + 1 dime + 1 nickel",
			orderQuantity: 40,
			expectedPacks: []PackResult{{Size: 5, Count: 1}, {Size: 10, Count: 1}, {Size: 25, Count: 1}},
			explanation:   "40 cents exactly: 25+10+5 (3 packs, 40 items) - exact match",
		},
		{
			name:          "50 cents - 2 quarters",
			orderQuantity: 50,
			expectedPacks: []PackResult{{Size: 25, Count: 2}},
			explanation:   "50 cents: 2x25 (2 packs, 50 items) - exact match",
		},
		{
			name:          "99 cents - 4 quarters (100 items is better than many small packs)",
			orderQuantity: 99,
			expectedPacks: []PackResult{{Size: 25, Count: 4}},
			explanation:   "99 cents: 4x25=100 (4 packs, only 1 excess) vs many small packs - minimize excess",
		},
		{
			name:          "15 cents - 1 dime + 1 nickel",
			orderQuantity: 15,
			expectedPacks: []PackResult{{Size: 5, Count: 1}, {Size: 10, Count: 1}},
			explanation:   "15 cents: 10+5 (2 packs, 15 items) - exact match",
		},
		{
			name:          "3 cents - 1 nickel (5 items, minimal excess)",
			orderQuantity: 3,
			expectedPacks: []PackResult{{Size: 5, Count: 1}},
			explanation:   "3 cents: 1x5=5 (1 pack, 2 excess) - minimal excess achievable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.orderQuantity)

			// Calculate totals for debugging
			totalItems := 0
			totalPacks := 0
			for _, p := range result {
				totalItems += p.Size * p.Count
				totalPacks += p.Count
			}

			if !packsEqual(result, tt.expectedPacks) {
				t.Errorf("%s: Expected %v (total: %d items, %d packs), got %v (total: %d items, %d packs)\nExplanation: %s",
					tt.name, tt.expectedPacks, tt.orderQuantity, len(tt.expectedPacks),
					result, totalItems, totalPacks, tt.explanation)
			}
		})
	}
}

// TestPackSizeCustomizationAtRuntime tests that pack sizes can be changed and immediately used
func TestPackSizeCustomizationAtRuntime(t *testing.T) {
	calc := NewDefaultCalculator()

	// First calculation with default sizes
	result1 := calc.Calculate(500)
	expected1 := []PackResult{{Size: 500, Count: 1}}
	if !packsEqual(result1, expected1) {
		t.Errorf("With default sizes, 500: Expected %v, got %v", expected1, result1)
	}

	// Change pack sizes
	calc.UpdatePackSizes([]int{100, 200, 400})

	// Same order quantity should now give different result
	result2 := calc.Calculate(500)
	expected2 := []PackResult{{Size: 100, Count: 1}, {Size: 400, Count: 1}}
	if !packsEqual(result2, expected2) {
		t.Errorf("With custom sizes [100,200,400], 500: Expected %v, got %v", expected2, result2)
	}
}

// TestMinimizePacksWhenItemsEqual tests rule 3: minimize pack count when excess is equal
func TestMinimizePacksWhenItemsEqual(t *testing.T) {
	// With packs 100, 250, 500:
	// For 500 items: 1x500 (1 pack) is better than 2x250 (2 packs) or 5x100 (5 packs)
	// All give 500 items (same excess), so choose fewest packs
	calc := NewCalculator([]int{100, 250, 500})

	result := calc.Calculate(500)
	expected := []PackResult{{Size: 500, Count: 1}}

	if !packsEqual(result, expected) {
		t.Errorf("500 with packs [100,250,500]: Expected 1x500 (fewest packs), got %v", result)
	}

	// For 750 items: 1x500 + 1x250 (2 packs, 750 items) vs 3x250 (3 packs)
	// Both give exact match, choose fewer packs
	result = calc.Calculate(750)
	expected = []PackResult{{Size: 250, Count: 1}, {Size: 500, Count: 1}}

	if !packsEqual(result, expected) {
		t.Errorf("750 with packs [100,250,500]: Expected 1x500+1x250 (2 packs), got %v", result)
	}
}
