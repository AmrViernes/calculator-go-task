package calculator

import (
	"math"
	"sort"
)

// PackConfig represents a single pack size (not really used, but here for completeness)
type PackConfig struct {
	Size int `json:"size"`
}

// PackResult is what we return - how many of each pack size to send
type PackResult struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

// Calculator holds the pack sizes and does the math
type Calculator struct {
	packSizes []int // Always kept sorted ascending
}

// NewCalculator creates a calculator with custom pack sizes
func NewCalculator(packSizes []int) *Calculator {
	// Make a copy so the caller can't modify our internal state
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Ints(sizes)

	return &Calculator{
		packSizes: sizes,
	}
}

// DefaultPackSizes returns the standard pack sizes we use
func DefaultPackSizes() []int {
	return []int{250, 500, 1000, 2000, 5000}
}

// NewDefaultCalculator creates a calculator with standard pack sizes
func NewDefaultCalculator() *Calculator {
	return NewCalculator(DefaultPackSizes())
}

// GetPackSizes returns a copy of current pack sizes
func (c *Calculator) GetPackSizes() []int {
	sizes := make([]int, len(c.packSizes))
	copy(sizes, c.packSizes)
	return sizes
}

// UpdatePackSizes changes the available pack sizes
func (c *Calculator) UpdatePackSizes(packSizes []int) {
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Ints(sizes)
	c.packSizes = sizes
}

// Calculate figures out the optimal pack combination for an order
// The rules are:
// 1. We must send at least as many items as ordered (can't under-ship)
// 2. Minimize excess items (don't send way more than needed)
// 3. When excess is equal, use fewer packs
func (c *Calculator) Calculate(orderQuantity int) []PackResult {
	// Edge cases: nothing to calculate
	if orderQuantity <= 0 || len(c.packSizes) == 0 {
		return []PackResult{}
	}

	// Find the best combination of packs
	result := c.findOptimalCombination(orderQuantity)

	// Convert the map to a clean result array
	var results []PackResult
	for _, size := range c.packSizes {
		if count := result[size]; count > 0 {
			results = append(results, PackResult{
				Size:  size,
				Count: count,
			})
		}
	}

	return results
}

// findOptimalCombination uses dynamic programming to find the best pack combo
// This is essentially the "unbounded knapsack" problem with a twist:
// we can exceed the target, but we want to minimize the excess
func (c *Calculator) findOptimalCombination(target int) map[int]int {
	// Zero order means no packs needed
	if target == 0 {
		return make(map[int]int)
	}

	maxSize := c.packSizes[len(c.packSizes)-1]
	// We might need to send more than the target, but not by more than one largest pack
	maxQuantity := target + maxSize

	// This struct tracks a potential combination
	type combination struct {
		packs      map[int]int // pack size -> count
		totalItems int         // total items in this combination
		packCount  int         // total number of packs
	}

	// memo[q] stores the best combination to get at least q items
	memo := make(map[int]*combination)

	// Base case: zero items needs zero packs
	memo[0] = &combination{
		packs:      make(map[int]int),
		totalItems: 0,
		packCount:  0,
	}

	// Build up solutions for each quantity
	for q := 1; q <= maxQuantity; q++ {
		best := (*combination)(nil)

		// Try adding each pack size
		for _, packSize := range c.packSizes {
			if packSize > q {
				continue
			}

			prev := q - packSize
			if memo[prev] == nil {
				continue
			}

			// Create a new combination by adding this pack
			newPacks := c.copyPacks(memo[prev].packs)
			newPacks[packSize]++

			newTotalItems := memo[prev].totalItems + packSize
			newPackCount := memo[prev].packCount + 1

			newComb := &combination{
				packs:      newPacks,
				totalItems: newTotalItems,
				packCount:  newPackCount,
			}

			// Pick the best option:
			// 1. Fewer total items (less excess)
			// 2. If items are equal, fewer packs
			if best == nil || newTotalItems < best.totalItems ||
				(newTotalItems == best.totalItems && newPackCount < best.packCount) {
				best = newComb
			}
		}

		if best != nil {
			memo[q] = best
		}
	}

	// Now find the best combination that meets or exceeds our target
	var bestComb *combination
	minExcess := math.MaxInt32
	minPacks := math.MaxInt32

	for q := target; q <= maxQuantity; q++ {
		if memo[q] == nil {
			continue
		}

		excess := memo[q].totalItems - target

		// Minimize excess first, then pack count
		if excess < minExcess || (excess == minExcess && memo[q].packCount < minPacks) {
			bestComb = memo[q]
			minExcess = excess
			minPacks = memo[q].packCount
		}
	}

	if bestComb == nil {
		return make(map[int]int)
	}

	return bestComb.packs
}

// copyPacks makes a deep copy of a pack map (so we don't share state)
func (c *Calculator) copyPacks(packs map[int]int) map[int]int {
	result := make(map[int]int)
	for k, v := range packs {
		result[k] = v
	}
	return result
}
