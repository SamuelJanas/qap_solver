package solvers

import (
	"fmt"
	"qap_solver/internal/qap"
)

type GreedySolver struct {
	MaxIterations  int
	RandomRestarts int
}

func NewGreedySolver(maxIterations, randomRestarts int) *GreedySolver {
	return &GreedySolver{
		MaxIterations:  maxIterations,
		RandomRestarts: randomRestarts,
	}
}

func (s *GreedySolver) Name() string {
	return "GreedySolver"
}

func (s *GreedySolver) Description() string {
	return fmt.Sprintf("Greedy search with max iterations: %d and random restarts: %d", s.MaxIterations, s.RandomRestarts)
}

func (s *GreedySolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for restart := 0; restart < s.RandomRestarts; restart++ {
		currentSolution := RandomSolution(instance.Size)
		currentFitness := qap.CalculateFitness(instance, currentSolution)

		for iter := 0; iter < s.MaxIterations; iter++ {
			improved := false
			for i := 0; i < instance.Size-1; i++ {
				for j := i + 1; j < instance.Size; j++ {
					newSolution := make([]int, instance.Size)
					copy(newSolution, currentSolution)
					newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
					newFitness := qap.CalculateFitness(instance, newSolution)

					if newFitness < currentFitness {
						copy(currentSolution, newSolution)
						currentFitness = newFitness
						improved = true
						break
					}
				}
				if improved {
					break
				}
			}
			if !improved {
				break
			}
		}
		if bestFitness == -1 || currentFitness < bestFitness {
			copy(bestSolution, currentSolution)
			bestFitness = currentFitness
		}
	}
	return SolverResult{Solution: bestSolution, Fitness: bestFitness}
}
