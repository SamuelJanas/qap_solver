package solvers

import (
	"fmt"
	"qap_solver/internal/qap"
)

// LocalSearchSolver implements a basic local search strategy with swaps
type LocalSearchSolver struct {
	MaxIterations  int
	MaxNonImproving int
	RandomRestarts int
}

// NewLocalSearchSolver creates a new local search solver
func NewLocalSearchSolver(maxIterations, maxNonImproving, randomRestarts int) *LocalSearchSolver {
	return &LocalSearchSolver{
		MaxIterations:   maxIterations,
		MaxNonImproving: maxNonImproving,
		RandomRestarts:  randomRestarts,
	}
}

func (s *LocalSearchSolver) Name() string {
	return "LocalSearch"
}

func (s *LocalSearchSolver) Description() string {
	return fmt.Sprintf("Local search with swap neighborhood (Max iterations: %d, Non-improving limit: %d, Random restarts: %d)",
		s.MaxIterations, s.MaxNonImproving, s.RandomRestarts)
}

// Solve runs the local search algorithm
func (s *LocalSearchSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for restart := 0; restart < s.RandomRestarts; restart++ {
		// Start with a random solution
		currentSolution := RandomSolution(instance.Size)
		currentFitness := qap.CalculateFitness(instance, currentSolution)

		nonImprovingCount := 0

		for iter := 0; iter < s.MaxIterations && nonImprovingCount < s.MaxNonImproving; iter++ {
			improved := false

			// Try all possible swaps to find improvement
			for i := 0; i < instance.Size-1; i++ {
				for j := i + 1; j < instance.Size; j++ {
					// Create a new solution by swapping positions i and j
					newSolution := make([]int, instance.Size)
					copy(newSolution, currentSolution)
					newSolution[i], newSolution[j] = newSolution[j], newSolution[i]

					// Calculate fitness of new solution
					newFitness := qap.CalculateFitness(instance, newSolution)

					// If it's better, accept it
					if newFitness < currentFitness {
						copy(currentSolution, newSolution)
						currentFitness = newFitness
						improved = true
						// Break out of inner loop to start again with the new solution
						break
					}
				}
				if improved {
					break
				}
			}

			// If no improvement was found in this iteration
			if !improved {
				nonImprovingCount++
			} else {
				nonImprovingCount = 0
			}
		}

		// Update best solution if this restart found a better one
		if bestFitness == -1 || currentFitness < bestFitness {
			copy(bestSolution, currentSolution)
			bestFitness = currentFitness
		}
	}

	return SolverResult{
		Solution: bestSolution,
		Fitness:  bestFitness,
	}
}
