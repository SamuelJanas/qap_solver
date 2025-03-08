package solvers

import (
    "math/rand"
)

// Generate a random solution (random permutation)
func RandomSolution(size int) []int {
    solution := make([]int, size)
    for i := range solution {
        solution[i] = i
    }
    rand.Shuffle(size, func(i, j int) { solution[i], solution[j] = solution[j], solution[i] })
    return solution
}
