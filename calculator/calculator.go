package calculator

import (
	"math"
	"sort"
)

type PackConfig struct {
	Size int `json:"size"`
}

type PackResult struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

type Calculator struct {
	packSizes []int
}

func NewCalculator(packSizes []int) *Calculator {
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Ints(sizes)

	return &Calculator{
		packSizes: sizes,
	}
}

func DefaultPackSizes() []int {
	return []int{250, 500, 1000, 2000, 5000}
}

func NewDefaultCalculator() *Calculator {
	return NewCalculator(DefaultPackSizes())
}

func (c *Calculator) GetPackSizes() []int {
	sizes := make([]int, len(c.packSizes))
	copy(sizes, c.packSizes)
	return sizes
}

func (c *Calculator) UpdatePackSizes(packSizes []int) {
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Ints(sizes)
	c.packSizes = sizes
}

func (c *Calculator) Calculate(orderQuantity int) []PackResult {
	if orderQuantity <= 0 || len(c.packSizes) == 0 {
		return []PackResult{}
	}

	result := c.findOptimalCombination(orderQuantity)

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

func (c *Calculator) findOptimalCombination(target int) map[int]int {
	if target == 0 {
		return make(map[int]int)
	}

	maxSize := c.packSizes[len(c.packSizes)-1]
	maxQuantity := target + maxSize

	type combination struct {
		packs      map[int]int
		totalItems int
		packCount  int
	}

	memo := make(map[int]*combination)

	memo[0] = &combination{
		packs:      make(map[int]int),
		totalItems: 0,
		packCount:  0,
	}

	for q := 1; q <= maxQuantity; q++ {
		best := (*combination)(nil)

		for _, packSize := range c.packSizes {
			if packSize > q {
				continue
			}

			prev := q - packSize
			if memo[prev] == nil {
				continue
			}

			newPacks := c.copyPacks(memo[prev].packs)
			newPacks[packSize]++

			newTotalItems := memo[prev].totalItems + packSize
			newPackCount := memo[prev].packCount + 1

			newComb := &combination{
				packs:      newPacks,
				totalItems: newTotalItems,
				packCount:  newPackCount,
			}

			if best == nil || newTotalItems < best.totalItems ||
				(newTotalItems == best.totalItems && newPackCount < best.packCount) {
				best = newComb
			}
		}

		if best != nil {
			memo[q] = best
		}
	}

	var bestComb *combination
	minExcess := math.MaxInt32
	minPacks := math.MaxInt32

	for q := target; q <= maxQuantity; q++ {
		if memo[q] == nil {
			continue
		}

		excess := memo[q].totalItems - target

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

func (c *Calculator) copyPacks(packs map[int]int) map[int]int {
	result := make(map[int]int)
	for k, v := range packs {
		result[k] = v
	}
	return result
}
