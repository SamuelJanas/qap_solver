package solvers

import (
	"fmt"
	"math/rand"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"time"
)

type RandomWalkSolver struct {
	MaxIterations  int
	RandomRestarts int
}

func NewRandomWalkSolver(maxIterations int) *RandomWalkSolver {
	return &RandomWalkSolver{
		MaxIterations: maxIterations,
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

func (s *RandomWalkSolver) SolveWithMetrics(
	instance *qap.QAPInstance,
	metricsCollector *metrics.MetricsCollector,
	instanceName string,
	runNumber int,
) SolverResult {
	startTime := time.Now()

	// Initial values for solution and fitness
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	currentSolution := RandomSolution(instance.Size)
	currentFitness := qap.CalculateFitness(instance, currentSolution)

	// Metrics counters
	totalSteps := 0
	totalEvaluations := 0
	totalSolutionsChecked := 0

	var initialSolution []int
	var initialFitness int

	// Record initial solution
	initialSolution = make([]int, len(currentSolution))
	copy(initialSolution, currentSolution)
	initialFitness = currentFitness

	// Start the random walk search
	for iter := 0; iter < s.MaxIterations; iter++ {
		// Randomly select two indices i and j
		i, j := rand.Intn(instance.Size), 1+rand.Intn(instance.Size-2)
		j = (i + j) % instance.Size

		// Generate a new solution by swapping i and j
		newSolution := make([]int, instance.Size)
		copy(newSolution, currentSolution)
		newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
		newFitness := qap.CalculateFitness(instance, newSolution)

		totalEvaluations++
		totalSolutionsChecked++

		// Accept the new solution
		copy(currentSolution, newSolution)
		currentFitness = newFitness

		// If the new solution is better, update the best solution
		if bestFitness == -1 || currentFitness < bestFitness {
			copy(bestSolution, currentSolution)
			bestFitness = currentFitness
		}

		totalSteps++
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	// Record metrics if the collector is provided
	if metricsCollector != nil {
		metricsCollector.AddRunMetrics(metrics.RunMetrics{
			InstanceName:     instanceName,
			SolverName:       s.Name(),
			Run:              runNumber,
			InitialFitness:   initialFitness,
			FinalFitness:     bestFitness,
			TimeElapsed:      elapsedTime,
			StepsCount:       totalSteps,
			EvaluationsCount: totalEvaluations,
			SolutionsChecked: totalSolutionsChecked,
			Solution:         bestSolution,
		})
	}

	// Return the result
	return SolverResult{
		Solution: bestSolution,
		Fitness:  bestFitness,
	}
}
