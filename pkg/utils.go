package pkg

import (
	"math/rand"
)

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomIntPair(min, max int) (int, int) {
	if max-min < 1 {
		panic("Range too small to generate two different numbers")
	}

	first := RandomInt(min, max)
	second := first

	// Faster than modulo for larger instances.
	// The infinite loop is inplausible
	for second == first {
		second = RandomInt(min, max)
	}

	return first, second
}

func ShuffleSlice(slice []int) {
	n := len(slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
