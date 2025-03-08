package solvers

import (
	"fmt"
	"qap_solver/internal/qap"
	"qap_solver/pkg"
)

type Baseline struct {
	Arg1 int
	Arg2 int
	Arg3 string
}

// NewBaseline creates a new random solver with 3 arguments
func NewBaseline(arg1, arg2 int, arg3 string) *Baseline {
	return &Baseline{
		Arg1: arg1,
		Arg2: arg2,
		Arg3: arg3,
	}
}

// Name specifies the name of the method
func (s *Baseline) Name() string {
	return "baseline"
}

// Description provides description of the method
func (s *Baseline) Description() string {
	return fmt.Sprint("Showcase of functionality")
}

// Solves is the function called to solve the instance
func (s *Baseline) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for i := 0; i < s.Arg1; i++ {
		solution := SomeSolutionFunction(instance.Size)
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

// helper functions start here

func SomeSolutionFunction(size int) []int {
	solution := make([]int, size)
	for i := range solution {
		solution[i] = i
	}
	pkg.ShuffleSlice(solution)
	return solution
}
