package solvers

import (
	"fmt"
	"math/rand"
	"qap_solver/internal/qap"
)

type RandomWalkSolver struct {
	MaxIterations  int
	RandomRestarts int
}

func NewRandomWalkSolver(maxIterations, randomRestarts int) *RandomWalkSolver {
	return &RandomWalkSolver{
		MaxIterations:  maxIterations,
		RandomRestarts: randomRestarts,
	}
}

func (s *RandomWalkSolver) Name() string {
	return "RandomWalkSolver"
}

func (s *RandomWalkSolver) Description() string {
	return fmt.Sprintf("Random walk search with max iterations: %d and random restarts: %d", s.MaxIterations, s.RandomRestarts)
}

func (s *RandomWalkSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for restart := 0; restart < s.RandomRestarts; restart++ {
		currentSolution := RandomSolution(instance.Size)
		currentFitness := qap.CalculateFitness(instance, currentSolution)

		for iter := 0; iter < s.MaxIterations; iter++ {
			i, j := rand.Intn(instance.Size), 1+rand.Intn(instance.Size-2)
			j = (i + j) % instance.Size

			newSolution := make([]int, instance.Size)
			copy(newSolution, currentSolution)
			newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
			newFitness := qap.CalculateFitness(instance, newSolution)

			if newFitness < currentFitness || rand.Float64() < 0.1 {
				copy(currentSolution, newSolution)
				currentFitness = newFitness
			}
		}

		if bestFitness == -1 || currentFitness < bestFitness {
			copy(bestSolution, currentSolution)
			bestFitness = currentFitness
		}
	}

	return SolverResult{Solution: bestSolution, Fitness: bestFitness}
}
