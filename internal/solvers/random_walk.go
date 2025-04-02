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

func NewRandomWalkSolver(maxIterations int) *RandomWalkSolver {
	return &RandomWalkSolver{
		MaxIterations:  maxIterations,
	}
}

func (s *RandomWalkSolver) Name() string {
	return "Random Walk"
}

func (s *RandomWalkSolver) Description() string {
	return fmt.Sprintf("Random walk search with max iterations: %d", s.MaxIterations)
}

func (s *RandomWalkSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	currentSolution := RandomSolution(instance.Size)
	currentFitness := qap.CalculateFitness(instance, currentSolution)

	for iter := 0; iter < s.MaxIterations; iter++ {
		i, j := rand.Intn(instance.Size), 1+rand.Intn(instance.Size-2)
		j = (i + j) % instance.Size

		newSolution := make([]int, instance.Size)
		copy(newSolution, currentSolution)
		newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
		newFitness := qap.CalculateFitness(instance, newSolution)

		copy(currentSolution, newSolution)
		currentFitness = newFitness

		if bestFitness == -1 || currentFitness < bestFitness {
			copy(bestSolution, currentSolution)
			bestFitness = currentFitness
		}
		
	}

	return SolverResult{Solution: bestSolution, Fitness: bestFitness}
}
