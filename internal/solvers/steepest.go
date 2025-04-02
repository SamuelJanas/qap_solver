package solvers

import (
	"fmt"
	"qap_solver/internal/qap"
)

type SteepestSolver struct {
	MaxIterations  int
	RandomRestarts int
}

func NewSteepestSolver(maxIterations, randomRestarts int) *SteepestSolver {
	return &SteepestSolver{
		MaxIterations:  maxIterations,
		RandomRestarts: randomRestarts,
	}
}

func (s *SteepestSolver) Name() string {
	return "SteepestSolver"
}

func (s *SteepestSolver) Description() string {
	return fmt.Sprintf("Steepest search with max iterations: %d and random restarts: %d", s.MaxIterations, s.RandomRestarts)
}

func (s *SteepestSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for restart := 0; restart < s.RandomRestarts; restart++ {
		currentSolution := RandomSolution(instance.Size)
		currentFitness := qap.CalculateFitness(instance, currentSolution)

		for iter := 0; iter < s.MaxIterations; iter++ {
			bestNeighbor := make([]int, instance.Size)
			copy(bestNeighbor, currentSolution)
			bestNeighborFitness := currentFitness

			for i := 0; i < instance.Size-1; i++ {
				for j := i + 1; j < instance.Size; j++ {
					newSolution := make([]int, instance.Size)
					copy(newSolution, currentSolution)
					newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
					newFitness := qap.CalculateFitness(instance, newSolution)

					if newFitness < bestNeighborFitness {
						copy(bestNeighbor, newSolution)
						bestNeighborFitness = newFitness
					}
				}
			}
			if bestNeighborFitness < currentFitness {
				copy(currentSolution, bestNeighbor)
				currentFitness = bestNeighborFitness
			} else {
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
