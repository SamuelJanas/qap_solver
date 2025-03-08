package solvers

import (
	"fmt"
	"qap_solver/internal/qap"
	"qap_solver/pkg"
)

type RandomSolver struct {
	Iterations int
}

// NewRandomSolver creates a new random solver with specified iterations
func NewRandomSolver(iterations int) *RandomSolver {
	return &RandomSolver{
		Iterations: iterations,
	}
}

func (s *RandomSolver) Name() string {
	return "Random"
}

func (s *RandomSolver) Description() string {
	return fmt.Sprintf("Random solution generator (%d iterations)", s.Iterations)
}

func (s *RandomSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for i := 0; i < s.Iterations; i++ {
		solution := RandomSolution(instance.Size)
		fitness := qap.CalculateFitness(instance, solution)

		if bestFitness == -1 || fitness < bestFitness {
			copy(bestSolution, solution)
			bestFitness = fitness
		}
	}

	return SolverResult{
		Solution: bestSolution,
		Fitness:  bestFitness,
	}
}

func RandomSolution(size int) []int {
	solution := make([]int, size)
	for i := range solution {
		solution[i] = i
	}
	pkg.ShuffleSlice(solution)
	return solution
}
